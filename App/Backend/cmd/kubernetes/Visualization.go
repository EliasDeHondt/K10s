package kubernetes

import (
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"sync"
	"time"
)

type Visualization struct {
	Cluster  *ClusterView
	Services []*ServiceView
	mu       sync.Mutex
}

func (v *Visualization) findOrCreateNodeView(nodeName, namespace string) *NodeView {
	for _, node := range v.Cluster.Nodes {
		if node.Name == nodeName && node.Namespace == namespace {
			return node
		}
	}

	newNode := &NodeView{
		Name:        nodeName,
		Namespace:   namespace,
		Deployments: []*DeploymentView{},
	}
	v.Cluster.Nodes = append(v.Cluster.Nodes, newNode)
	return newNode
}

func (v *Visualization) findOrCreateServiceView(serviceName string) *ServiceView {
	for _, service := range v.Services {
		if service.Name == serviceName {
			return service
		}
	}

	// If not found, create a new ServiceView
	newService := &ServiceView{
		Name:          serviceName,
		Deployments:   []*DeploymentView{},
		LoadBalancers: []*LoadBalancer{},
	}
	v.Services = append(v.Services, newService)
	return newService
}

func matchLabels(selector, labels map[string]string) bool {
	for key, value := range selector {
		if labels[key] != value {
			return false
		}
	}
	return true
}

func (v *Visualization) AddNode(node *v1.Node, client IClient) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.Cluster.Nodes = append(v.Cluster.Nodes, NewNodeView(node, client))
}

func (v *Visualization) DeleteNode(node *v1.Node) {
	v.mu.Lock()
	defer v.mu.Unlock()

	for i, nodeView := range v.Cluster.Nodes {
		if nodeView.Name == node.Name {
			v.Cluster.Nodes = append(v.Cluster.Nodes[:i], v.Cluster.Nodes[i+1:]...)
			break
		}
	}
}

func (v *Visualization) AddService(service *v1.Service, client IClient) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.Services = append(v.Services, NewServiceView(service, client))
}

func (v *Visualization) DeleteService(service *v1.Service) {
	v.mu.Lock()
	defer v.mu.Unlock()

	for i, serviceView := range v.Services {
		if serviceView.Name == service.Name {
			v.Services = append(v.Services[:i], v.Services[i+1:]...)
			break
		}
	}
}

func removeDeploymentFromList(deployments []*DeploymentView, deploymentName string) []*DeploymentView {
	result := make([]*DeploymentView, 0)
	for _, deployment := range deployments {
		if deployment.Name != deploymentName {
			result = append(result, deployment)
		}
	}
	return result
}

func (v *Visualization) AddDeployment(deployment *appsv1.Deployment, client IClient) {
	v.mu.Lock()
	defer v.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newDeployment := NewDeploymentView(deployment)

	podList, err := client.GetPods(deployment.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		log.Printf("Error getting Pods for Deployment %s: %v", deployment.Name, err)
		return
	}

	for _, pod := range podList.Items {
		nodeName := pod.Spec.NodeName
		namespace := pod.Namespace

		nodeView := v.findOrCreateNodeView(nodeName, namespace)
		nodeView.Deployments = append(nodeView.Deployments, newDeployment)
	}

	serviceList, err := client.GetServices(deployment.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("Error getting Services for Deployment %s: %v", deployment.Name, err)
		return
	}

	for _, service := range serviceList.Items {
		if matchLabels(service.Spec.Selector, deployment.Spec.Template.Labels) {
			serviceView := v.findOrCreateServiceView(service.Name)
			serviceView.Deployments = append(serviceView.Deployments, newDeployment)
		}
	}

}

func (v *Visualization) DeleteDeployment(deployment *appsv1.Deployment) {
	v.mu.Lock()
	defer v.mu.Unlock()

	deploymentName := deployment.Name

	for _, node := range v.Cluster.Nodes {
		node.Deployments = removeDeploymentFromList(node.Deployments, deploymentName)
	}

	for _, service := range v.Services {
		service.Deployments = removeDeploymentFromList(service.Deployments, deploymentName)
	}
}

type ClusterView struct {
	Name  string
	Nodes []*NodeView
}

type NodeView struct {
	Name        string
	Namespace   string
	Deployments []*DeploymentView
}

type LoadBalancer struct {
	HostName string
	IP       string
	Services []*ServiceView
}

type ServiceView struct {
	Name          string
	Deployments   []*DeploymentView
	LoadBalancers []*LoadBalancer
}

type DeploymentView struct {
	Name string
}

func VisualizeCluster(client IClient) *Visualization {

	return &Visualization{
		Cluster:  NewClusterView(client),
		Services: getAllServices(client),
	}
}

func NewClusterView(client IClient) *ClusterView {

	nodes, err := createNodeViews(client)
	if err != nil {
		log.Fatal(err)
	}

	return &ClusterView{
		Name:  "Cluster",
		Nodes: nodes,
	}
}

func NewNodeView(node *v1.Node, client IClient) *NodeView {

	return &NodeView{
		Name:        node.Name,
		Namespace:   node.Namespace,
		Deployments: linkDeploymentsToNodes(client, node.Name),
	}
}

func NewServiceView(service *v1.Service, client IClient) *ServiceView {

	if len(service.Spec.Selector) == 0 {
		return &ServiceView{
			Name:          service.Name,
			Deployments:   make([]*DeploymentView, 0),
			LoadBalancers: make([]*LoadBalancer, 0),
		}
	}

	deployments, err := getDeploymentsForService(service.Namespace, service.Name, client)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	loadBalancers, err := getLoadBalancersForService(service)
	if err != nil {
		log.Fatal(err)
	}

	return &ServiceView{
		Name:          service.Name,
		Deployments:   deployments,
		LoadBalancers: loadBalancers,
	}
}

func NewLoadBalancer(ingress *v1.LoadBalancerIngress) *LoadBalancer {
	return &LoadBalancer{
		HostName: ingress.Hostname,
		IP:       ingress.IP,
	}
}

func getLoadBalancersForService(service *v1.Service) ([]*LoadBalancer, error) {
	loadBalancers := make([]*LoadBalancer, 0)
	balancers := service.Status.LoadBalancer.Ingress

	if len(balancers) > 0 {
		for _, ingress := range balancers {
			loadBalancers = append(loadBalancers, NewLoadBalancer(&ingress))
		}
	}

	return loadBalancers, nil
}

func getAllServices(client IClient) []*ServiceView {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services, err := client.GetServices("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list services: %v", err)
	}

	serviceList := make([]*ServiceView, 0, len(services.Items))
	for _, service := range services.Items {
		serviceList = append(serviceList, NewServiceView(&service, client))
	}

	return serviceList
}

func getDeploymentsForService(namespace, serviceName string, client IClient) ([]*DeploymentView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service, err := client.GetServices(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	serviceSelector := service.Spec.Selector
	if len(serviceSelector) == 0 {
		return []*DeploymentView{}, nil
	}
	serviceLabelSelector := labels.Set(serviceSelector).AsSelector().String()

	pods, err := client.GetPods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: serviceLabelSelector,
	})

	if err != nil {
		return nil, err
	}

	podLabels := make(map[string]map[string]string)
	for _, pod := range pods.Items {
		podLabels[pod.Name] = pod.Labels
	}

	deployments, err := client.GetDeployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	deploymentViews := transformDeployments(&deployments.Items, podLabels)

	return deploymentViews, nil
}

func transformDeployments(deployments *[]appsv1.Deployment, podLabels map[string]map[string]string) []*DeploymentView {
	deploymentViews := make([]*DeploymentView, 0)
	for _, deployment := range *deployments {
		deploySelector := deployment.Spec.Selector.MatchLabels
		matches := false

		for _, lbl := range podLabels {
			match := true
			for key, value := range deploySelector {
				if podValue, exists := lbl[key]; !exists || podValue != value {
					match = false
					break
				}
			}
			if match {
				matches = true
				break
			}
		}

		if matches {
			deploymentViews = append(deploymentViews, NewDeploymentView(&deployment))
		}
	}

	return deploymentViews
}

func NewDeploymentView(deployment *appsv1.Deployment) *DeploymentView {
	return &DeploymentView{
		Name: deployment.Name,
	}
}

func createNodeViews(client IClient) ([]*NodeView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodes, err := client.GetNodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodeViews := make([]*NodeView, len(nodes.Items))

	for i, node := range nodes.Items {
		nodeViews[i] = NewNodeView(&node, client)
	}
	return nodeViews, nil
}

func linkDeploymentsToNodes(client IClient, name string) []*DeploymentView {
	deployments := make([]*DeploymentView, 0)

	deploymentList, err := client.GetDeployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, deployment := range deploymentList.Items {

		pods, err := getPodsForDeployment(client, &deployment)
		if err != nil {
			panic(err)
		}

		for _, pod := range pods {
			nodeName := pod.Spec.NodeName
			if nodeName == name {
				deployments = append(deployments, NewDeploymentView(&deployment))
			}
		}

	}

	return deployments
}

func getPodsForDeployment(client IClient, deployment *appsv1.Deployment) ([]*v1.Pod, error) {
	namespace := deployment.GetNamespace()
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	podList, err := client.GetPods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}

	pods := make([]*v1.Pod, 0, len(podList.Items))
	for _, pod := range podList.Items {
		pods = append(pods, &pod)
	}

	return pods, nil
}
