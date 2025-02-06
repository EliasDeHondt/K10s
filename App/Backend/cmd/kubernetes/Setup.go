package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"time"
)

func TestFakeClient() *fake.Clientset {
	clientset := fake.NewClientset()

	// Create fake nodes with labels and IPs
	node1 := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-1",
			Labels: map[string]string{
				"kubernetes.io/role": "worker",
				"version":            "v1.25.0",
			},
			CreationTimestamp: metav1.NewTime(time.Now().Add(-48 * time.Hour)), // Simulated 2-day-old node
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: corev1.ConditionTrue,
				},
			},
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "192.168.1.1"},
			},
		},
	}

	node2 := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-2",
			Labels: map[string]string{
				"kubernetes.io/role": "master",
				"version":            "v1.25.0",
			},
			CreationTimestamp: metav1.NewTime(time.Now().Add(-24 * time.Hour)), // 1-day-old node
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: corev1.ConditionFalse,
				},
			},
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "192.168.1.2"},
			},
		},
	}

	clientset.CoreV1().Nodes().Create(context.TODO(), node1, metav1.CreateOptions{})
	clientset.CoreV1().Nodes().Create(context.TODO(), node2, metav1.CreateOptions{})

	return clientset
}
