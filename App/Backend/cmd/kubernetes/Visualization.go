package kubernetes

import (
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"time"
)

type Visualization struct {
	Cluster  *ClusterView
	Services []*ServiceView
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
			Name:        service.Name,
			Deployments: make([]*DeploymentView, 0),
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
