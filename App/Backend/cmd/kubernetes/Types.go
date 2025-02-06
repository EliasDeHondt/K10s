package kubernetes

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"time"
)

type Node struct {
	ClusterName string
	Name        string
	Status      string
	IP          string
	Role        string
	Version     string
	NodeAge     string
}

func NewNode(node v1.Node) Node {
	return Node{
		ClusterName: node.Namespace,
		Name:        node.Name,
		Status:      isNodeOnline(&node),
		IP:          node.Status.Addresses[0].Address,
		Role:        node.Labels["kubernetes.io/role"],
		Version:     node.ObjectMeta.Labels["version"],
		NodeAge:     calculateNodeAge(&node),
	}
}

func calculateNodeAge(node *v1.Node) string {
	age := time.Since(node.CreationTimestamp.Time)

	return fmt.Sprintf("%d:%d:%d", int(age.Hours()/24), int(age.Hours())%24, int(age.Minutes())%60)
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
