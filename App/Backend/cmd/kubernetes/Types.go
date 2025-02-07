package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type Node struct {
	Name       string
	Status     string
	Role       string
	Version    string
	PodsAmount int
	NodeAge    string
	IP         string
}

func NewNode(node v1.Node, clientset *kubernetes.Clientset) Node {
	return Node{
		Name:       node.Name,
		Status:     isNodeOnline(&node),
		Role:       node.Labels["kubernetes.io/role"],
		Version:    node.ObjectMeta.Labels["version"],
		PodsAmount: getPodsInNode(node.Name, clientset),
		NodeAge:    calculateNodeAge(&node),
		IP:         node.Status.Addresses[0].Address,
	}
}

func calculateNodeAge(node *v1.Node) string {
	age := time.Since(node.CreationTimestamp.Time)

	return fmt.Sprintf("%d:%d:%d", int(age.Hours()/24), int(age.Hours())%24, int(age.Minutes())%60)
}

func getPodsInNode(nodeName string, clientset *kubernetes.Clientset) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return 0
	}

	return len(pods.Items)
}

func isNodeOnline(node *v1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			switch condition.Status {
			case v1.ConditionTrue:
				return "ONLINE ✅"
			case v1.ConditionFalse:
				return "OFFLINE ❌"
			case v1.ConditionUnknown:
				return "UNKNOWN ⚠️"
			}
		}
	}
	return "NO STATUS ❓"
}

type Pod struct {
}

type ClusterStructure struct {
	Name  string
	Nodes []NodeTree
}

type NodeTree struct {
	Node Node
	Pods []Pod
}

type PodTree struct {
	Pod Pod
}
