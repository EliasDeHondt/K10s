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

type ClusterView struct {
	Name  string
	Nodes []*NodeView
}

type NodeView struct {
	Name        string
	Namespace   string
	Services    []*ServiceView
	Deployments []*DeploymentView
}

type LoadBalancer struct {
	HostName string
	IP       string
}

type ServiceView struct {
	Name          string
	LoadBalancers []*LoadBalancer
	Deployments   []*DeploymentView
}

type DeploymentView struct {
	Name string
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

	serviceViews, err := getServicesOnNode(node.Name, client)
	if err != nil {
		log.Fatal(err)
	}

	return &NodeView{
		Name:      node.Name,
		Namespace: node.Namespace,
		Services:  serviceViews,
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
