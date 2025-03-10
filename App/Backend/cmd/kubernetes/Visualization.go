package kubernetes

import (
	"fmt"
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

	deployments, err := getDeploymentsOnNode(node.Name, client)
	if err != nil {
		log.Fatal(err)
	}

	return &NodeView{
		Name:        node.Name,
		Namespace:   node.Namespace,
		Deployments: deployments,
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

	balancers := service.Status.LoadBalancer.Ingress
	loadBalancers := make([]*LoadBalancer, 0)

	if len(balancers) > 0 {

		loadBalancers = make([]*LoadBalancer, len(balancers))

		for i, ingress := range balancers {
			loadBalancers[i] = NewLoadBalancer(&ingress)
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
		fmt.Println("Service has no selectors (e.g., headless service).")
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
	deploymentViews := make([]*DeploymentView, len(*deployments))
	for i, deployment := range *deployments {
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
			deploymentViews[i] = NewDeploymentView(&deployment)
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

func getDeploymentsOnNode(nodeName string, client IClient) ([]*DeploymentView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := client.GetPods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return nil, err
	}

	replicaSetMap := make(map[string]string)
	deploymentSet := make(map[string]appsv1.Deployment)

	replicaSets, err := client.GetReplicaSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, rs := range replicaSets.Items {
		for _, owner := range rs.OwnerReferences {
			if owner.Kind == "Deployment" {
				replicaSetMap[rs.Name] = owner.Name
			}
		}
	}

	for _, pod := range pods.Items {
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "ReplicaSet" {
				if deploymentName, exists := replicaSetMap[owner.Name]; exists {
					deployment, err := client.GetDeployments(pod.Namespace).Get(ctx, deploymentName, metav1.GetOptions{})
					if err != nil {
						log.Printf("Failed to get Deployment %s: %v", deploymentName, err)
						continue
					}
					deploymentSet[deployment.Name] = *deployment
				}
			}
		}
	}

	deploymentViews := make([]*DeploymentView, 0)
	for _, deployment := range deploymentSet {
		deploymentViews = append(deploymentViews, NewDeploymentView(&deployment))
	}

	return deploymentViews, nil

}
