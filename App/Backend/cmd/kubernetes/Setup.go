package kubernetes

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"time"
)

func TestFakeClient() IClient {
	var clientset IClient = &FakeClient{fake.NewClientset()}

	nodes := []*corev1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-1",
				Labels: map[string]string{
					"kubernetes.io/role": "worker",
					"version":            "v1.25.0",
				},
				CreationTimestamp: metav1.NewTime(time.Now().Add(-48 * time.Hour)),
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{Type: corev1.NodeReady, Status: corev1.ConditionTrue},
				},
				Addresses: []corev1.NodeAddress{
					{Type: corev1.NodeInternalIP, Address: "192.168.1.1"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-2",
				Labels: map[string]string{
					"kubernetes.io/role": "master",
					"version":            "v1.25.0",
				},
				CreationTimestamp: metav1.NewTime(time.Now().Add(-24 * time.Hour)),
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{Type: corev1.NodeReady, Status: corev1.ConditionFalse},
				},
				Addresses: []corev1.NodeAddress{
					{Type: corev1.NodeInternalIP, Address: "192.168.1.2"},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-3",
				Labels: map[string]string{
					"kubernetes.io/role": "master",
					"version":            "v1.25.0",
				},
				CreationTimestamp: metav1.NewTime(time.Now().Add(-48 * time.Hour)),
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{Type: corev1.NodeReady, Status: corev1.ConditionTrue},
				},
				Addresses: []corev1.NodeAddress{
					{Type: corev1.NodeInternalIP, Address: "192.168.1.3"},
				},
			},
		},
	}
	for _, node := range nodes {
		clientset.GetNodes().Create(context.TODO(), node, metav1.CreateOptions{})
	}

	pods := []*corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-1",
				Namespace: "default",
				Labels:    map[string]string{"app": "test-app"},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				PodIP: "10.0.0.1",
				ContainerStatuses: []corev1.ContainerStatus{
					{RestartCount: 2},
				},
			},
			Spec: corev1.PodSpec{NodeName: "node-1"},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-2",
				Namespace: "default",
				Labels:    map[string]string{"app": "test-app"},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodPending,
				PodIP: "10.0.0.2",
				ContainerStatuses: []corev1.ContainerStatus{
					{RestartCount: 1},
				},
			},
			Spec: corev1.PodSpec{NodeName: "node-2"},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-3",
				Namespace: "test",
				Labels:    map[string]string{"app": "test-app"},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodPending,
				PodIP: "10.0.0.3",
				ContainerStatuses: []corev1.ContainerStatus{
					{RestartCount: 1},
				},
			},
			Spec: corev1.PodSpec{NodeName: "node-3"},
		},
	}
	for _, pod := range pods {
		namespace := pod.GetNamespace()
		clientset.GetPods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	}

	services := []*corev1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "service-1", Namespace: "default"},
			Spec: corev1.ServiceSpec{
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "10.100.1.1",
				Selector:  map[string]string{"app": "test-app"},
			},
		},
	}
	for _, service := range services {
		clientset.GetServices("default").Create(context.TODO(), service, metav1.CreateOptions{})
	}

	var replicas int32 = 3

	deployments := []*appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "deployment-1", Namespace: "default"},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test-app"}},
			},
			Status: appsv1.DeploymentStatus{
				Replicas:          3,
				ReadyReplicas:     2,
				UpdatedReplicas:   3,
				AvailableReplicas: 2,
			},
		},
	}
	for _, deployment := range deployments {
		clientset.GetDeployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	}

	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "secret-1", Namespace: "default"},
			Type:       corev1.SecretTypeOpaque,
			Data:       map[string][]byte{"password": []byte("supersecret")},
		},
	}
	for _, secret := range secrets {
		clientset.GetSecrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
	}

	configMaps := []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "configmap-1", Namespace: "default"},
			Data:       map[string]string{"config": "value"},
		},
	}
	for _, configMap := range configMaps {
		clientset.GetConfigMaps("default").Create(context.TODO(), configMap, metav1.CreateOptions{})
	}

	return clientset
}
