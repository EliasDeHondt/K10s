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
	Cluster       *ClusterView
	LoadBalancers []*LoadBalancer
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
	Name        string
	Deployments []*DeploymentView
}

type DeploymentView struct {
	Name string
}

func VisualizeCluster(client IClient) *Visualization {

	return &Visualization{
		Cluster:       NewClusterView(client),
		LoadBalancers: GetAllLoadBalancers(client),
	}
}

func GetAllLoadBalancers(client IClient) []*LoadBalancer {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceList, err := client.GetServices("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list services: %v", err)
	}

	lbMap := buildLoadBalancerMap(serviceList, client)

	return mapToSlice(lbMap)
}

func buildLoadBalancerMap(svcList *v1.ServiceList, client IClient) map[string]*LoadBalancer {
	lbMap := make(map[string]*LoadBalancer)

	for _, service := range svcList.Items {
		if service.Spec.Type == v1.ServiceTypeLoadBalancer {
			for _, ingress := range service.Status.LoadBalancer.Ingress {
				key := ingress.IP
				if key == "" {
					key = ingress.Hostname
				}
				if key == "" {
					continue
				}

				if _, exists := lbMap[key]; !exists {
					lbMap[key] = NewLoadBalancer(&ingress, client)
				}
			}
		}
	}

	return lbMap
}

func mapToSlice(lbMap map[string]*LoadBalancer) []*LoadBalancer {
	loadBalancers := make([]*LoadBalancer, 0, len(lbMap))
	for _, lb := range lbMap {
		loadBalancers = append(loadBalancers, lb)
	}
	return loadBalancers
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

	return &ServiceView{
		Name:        service.Name,
		Deployments: deployments,
	}
}

func NewLoadBalancer(ingress *v1.LoadBalancerIngress, client IClient) *LoadBalancer {
	return &LoadBalancer{
		HostName: ingress.Hostname,
		IP:       ingress.IP,
		Services: getServicesForLoadBalancer(client, ingress),
	}
}

func getServicesForLoadBalancer(client IClient, ingress *v1.LoadBalancerIngress) []*ServiceView {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var services []*ServiceView

	svcList, err := client.GetServices("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list services: %v", err)
	}

	for _, svc := range svcList.Items {
		if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
			for _, lbIngress := range svc.Status.LoadBalancer.Ingress {
				if lbIngress.IP == ingress.IP || lbIngress.Hostname == ingress.Hostname {
					services = append(services, NewServiceView(&svc, client))
				}
			}
		}
	}

	return services
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

func getServicesOnNode(nodeName string, client IClient) ([]*ServiceView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := client.GetPods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})

	if err != nil {
		return nil, err
	}

	podLabels := make(map[string]map[string]string)
	for _, pod := range pods.Items {
		podLabels[pod.Name] = pod.Labels
	}

	services, err := client.GetServices("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, svc := range services.Items {
		selector := svc.Spec.Selector
		if len(selector) == 0 {
			continue
		}

		if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
			continue
		}

		for _, lbls := range podLabels {
			matches := true
			for key, value := range selector {
				if podValue, exists := lbls[key]; !exists || podValue != value {
					matches = false
					break
				}
			}

			if matches {
				fmt.Printf("- %s (Namespace: %s)\n", svc.Name, svc.Namespace)
				break
			}
		}
	}

	serviceViews := transformServices(&services.Items, podLabels, client)
	return serviceViews, nil
}

func transformServices(services *[]v1.Service, podLbls map[string]map[string]string, client IClient) []*ServiceView {

	serviceViews := make([]*ServiceView, len(*services))

	for i, svc := range *services {
		selector := svc.Spec.Selector
		if len(selector) == 0 {
			continue
		}

		for _, lbls := range podLbls {
			matches := true
			for key, value := range selector {
				if podValue, exists := lbls[key]; !exists || podValue != value {
					matches = false
					break
				}
			}

			if matches {
				serviceViews[i] = NewServiceView(&svc, client)
				break
			}
		}
	}

	return serviceViews
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
