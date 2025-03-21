/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package client

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func NewNode(node v1.Node, clientset IClient) Node {
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

func getPodsInNode(nodeName string, clientset IClient) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pods, err := (clientset).GetPods("").List(ctx, metav1.ListOptions{
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
