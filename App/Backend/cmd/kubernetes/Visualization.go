package kubernetes

import (
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
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

func (v *Visualization) findOrCreateServiceView(serviceName string, namespace string) *ServiceView {
	for _, service := range v.Services {
		if service.Name == serviceName {
			return service
		}
	}

	newService := &ServiceView{
		Name:          serviceName,
		Namespace:     namespace,
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
			serviceView := v.findOrCreateServiceView(service.Name, service.Namespace)
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
	Name            string
	Nodes           []*NodeView
	Endpoints       v1.Endpoints
	ControlPlaneURL string
	APIVersion      string
	Timeout         time.Duration
	QPS             float32
	Burst           int
}

type NodeView struct {
	Name         string
	Namespace    string
	Deployments  []*DeploymentView
	NodeInfo     v1.NodeSystemInfo
	NodeStatus   []v1.NodeCondition
	NodeAddress  []v1.NodeAddress
	ResourceList v1.ResourceList
}

type LoadBalancer struct {
	HostName  string
	IP        string
	Namespace string
}

type ServiceView struct {
	Name          string
	Namespace     string
	Deployments   []*DeploymentView
	LoadBalancers []*LoadBalancer
	Type          string
	ClusterIP     string
	ExternalIPs   []string
	ServiceStatus []metav1.Condition
}

type DeploymentView struct {
	Name              string
	Namespace         string
	Replicas          int32
	ReadyReplicas     int32
	UpdatedReplicas   int32
	AvailableReplicas int32
}

func VisualizeCluster(client IClient, config *rest.Config) *Visualization {
	return &Visualization{
		Cluster:  NewClusterView(client, config),
		Services: getAllServices(client),
	}
}

func NewClusterView(client IClient, config *rest.Config) *ClusterView {
	if config == nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Printf("Error inferring in-cluster config: %v", err)
			return &ClusterView{
				Name: "Cluster",
			}
		}
		config = inClusterConfig
	}

	controlPlaneURL := config.Host
	if controlPlaneURL == "" {
		log.Printf("Warning: No control plane URL found in config")
		return &ClusterView{
			Name: "Cluster",
		}
	}

	nodes, err := createNodeViews(client)
	if err != nil {
		log.Fatal(err)
	}

	// set to default mock values for fake client
	qps := config.QPS
	if qps == 0 {
		qps = 5
	}
	burst := config.Burst
	if burst == 0 {
		burst = 10
	}
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &ClusterView{
		Name:            "Cluster",
		ControlPlaneURL: controlPlaneURL,
		APIVersion:      config.APIPath,
		Timeout:         timeout,
		QPS:             qps,
		Burst:           burst,
		Nodes:           nodes,
	}
}

func NewNodeView(node *v1.Node, client IClient) *NodeView {

	return &NodeView{
		Name:         node.Name,
		Namespace:    node.Namespace,
		Deployments:  linkDeploymentsToNodes(client, node.Name),
		NodeInfo:     node.Status.NodeInfo,
		NodeStatus:   node.Status.Conditions,
		NodeAddress:  node.Status.Addresses,
		ResourceList: node.Status.Capacity,
	}
}

func NewServiceView(service *v1.Service, client IClient) *ServiceView {

	if len(service.Spec.Selector) == 0 {
		return &ServiceView{
			Name:          service.Name,
			Deployments:   make([]*DeploymentView, 0),
			LoadBalancers: make([]*LoadBalancer, 0),
			Namespace:     service.Namespace,
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
		Namespace:     service.Namespace,
		Type:          string(service.Spec.Type),
		ClusterIP:     service.Spec.ClusterIP,
		ExternalIPs:   service.Spec.ExternalIPs,
		ServiceStatus: service.Status.Conditions,
		Deployments:   deployments,
		LoadBalancers: loadBalancers,
	}
}

func NewLoadBalancer(ingress *v1.LoadBalancerIngress, namespace string) *LoadBalancer {
	return &LoadBalancer{
		HostName:  ingress.Hostname,
		IP:        ingress.IP,
		Namespace: namespace,
	}
}

func getLoadBalancersForService(service *v1.Service) ([]*LoadBalancer, error) {
	loadBalancers := make([]*LoadBalancer, 0)
	balancers := service.Status.LoadBalancer.Ingress

	if len(balancers) > 0 {
		for _, ingress := range balancers {
			loadBalancers = append(loadBalancers, NewLoadBalancer(&ingress, service.Namespace))
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
		Name:              deployment.Name,
		Namespace:         deployment.Namespace,
		Replicas:          deployment.Status.Replicas,
		ReadyReplicas:     deployment.Status.ReadyReplicas,
		UpdatedReplicas:   deployment.Status.UpdatedReplicas,
		AvailableReplicas: deployment.Status.AvailableReplicas,
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

func (v *Visualization) FilterByNamespace(namespace string) *Visualization {
	v.mu.Lock()
	defer v.mu.Unlock()

	filtered := &Visualization{
		Cluster: &ClusterView{
			Name:  v.Cluster.Name,
			Nodes: make([]*NodeView, 0),
		},
		Services: make([]*ServiceView, 0),
	}

	for _, node := range v.Cluster.Nodes {
		nodeCopy := node.DeepCopy()

		filteredDeployments := make([]*DeploymentView, 0)
		for _, deployment := range nodeCopy.Deployments {
			if deployment.Namespace == namespace {
				filteredDeployments = append(filteredDeployments, deployment)
			}
		}
		nodeCopy.Deployments = filteredDeployments

		if len(nodeCopy.Deployments) > 0 {
			filtered.Cluster.Nodes = append(filtered.Cluster.Nodes, nodeCopy)
		}
	}

	for _, service := range v.Services {
		if service.Namespace == namespace {
			serviceCopy := service.DeepCopy()

			filteredDeployments := make([]*DeploymentView, 0)
			for _, deployment := range serviceCopy.Deployments {
				if deployment.Namespace == namespace {
					filteredDeployments = append(filteredDeployments, deployment)
				}
			}
			serviceCopy.Deployments = filteredDeployments

			filteredLoadBalancers := make([]*LoadBalancer, 0)
			for _, lb := range serviceCopy.LoadBalancers {
				if lb.Namespace == namespace {
					filteredLoadBalancers = append(filteredLoadBalancers, lb)
				}
			}
			serviceCopy.LoadBalancers = filteredLoadBalancers

			filtered.Services = append(filtered.Services, serviceCopy)
		}
	}

	return filtered
}

func (n *NodeView) DeepCopy() *NodeView {

	view := &NodeView{
		Name:         n.Name,
		Namespace:    n.Namespace,
		Deployments:  make([]*DeploymentView, len(n.Deployments)),
		NodeInfo:     n.NodeInfo,
		NodeStatus:   n.NodeStatus,
		NodeAddress:  n.NodeAddress,
		ResourceList: n.ResourceList,
	}

	for i, deployment := range n.Deployments {
		view.Deployments[i] = &DeploymentView{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		}
	}

	return view
}

func (s *ServiceView) DeepCopy() *ServiceView {
	view := &ServiceView{
		Name:          s.Name,
		Namespace:     s.Namespace,
		Deployments:   make([]*DeploymentView, len(s.Deployments)),
		LoadBalancers: make([]*LoadBalancer, len(s.LoadBalancers)),
		Type:          s.Type,
		ClusterIP:     s.ClusterIP,
		ExternalIPs:   s.ExternalIPs,
		ServiceStatus: s.ServiceStatus,
	}

	for i, deployment := range s.Deployments {
		view.Deployments[i] = &DeploymentView{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		}
	}

	for i, lb := range s.LoadBalancers {
		view.LoadBalancers[i] = &LoadBalancer{
			HostName:  lb.HostName,
			IP:        lb.IP,
			Namespace: lb.Namespace,
		}
	}

	return view
}
