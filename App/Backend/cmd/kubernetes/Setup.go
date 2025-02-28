/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package kubernetes

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/metrics"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"time"
)

func GetClients() IClient {
	config, err := rest.InClusterConfig()

	if err != nil {
		client := TestFakeClient()
		return client
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		client := TestFakeClient()
		return client
	}

	mc, err := metricsv.NewForConfig(config)
	if err != nil {
		client := TestFakeClient()
		return client
	}

	var client IClient = &Client{c, mc}
	return client
}

func TestFakeClient() IClient {

	fakeMetricsClient := &FakeMetricsClient{
		NodeMetrics: map[string]*metrics.NodeMetrics{
			"node-1": {
				Usage: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:     resource.MustParse("500m"),
					corev1.ResourceMemory:  resource.MustParse("2Gi"),
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			"node-2": {
				Usage: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:     resource.MustParse("300m"),
					corev1.ResourceMemory:  resource.MustParse("1Gi"),
					corev1.ResourceStorage: resource.MustParse("2Gi"),
				},
			},
			"node-3": {
				Usage: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:     resource.MustParse("1000m"),
					corev1.ResourceMemory:  resource.MustParse("4Gi"),
					corev1.ResourceStorage: resource.MustParse("3Gi"),
				},
			},
		},
		PodMetrics: map[string]*metrics.PodMetrics{
			"default/pod-1": {
				Containers: []metrics.ContainerMetrics{
					{
						Name: "container-1",
						Usage: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("512Mi"),
						},
					},
				},
			},
			"default/pod-2": {
				Containers: []metrics.ContainerMetrics{
					{
						Name: "container-1",
						Usage: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			},
		},
	}

	var clientset IClient = &FakeClient{fake.NewClientset(), fakeMetricsClient}

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
				Capacity: corev1.ResourceList{
					corev1.ResourceCPU:     resource.MustParse("1000m"),
					corev1.ResourceMemory:  resource.MustParse("8Gi"),
					corev1.ResourceStorage: resource.MustParse("1Gi"),
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
				Capacity: corev1.ResourceList{
					corev1.ResourceCPU:     resource.MustParse("1000m"),
					corev1.ResourceMemory:  resource.MustParse("8Gi"),
					corev1.ResourceStorage: resource.MustParse("3Gi"),
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
				Capacity: corev1.ResourceList{
					corev1.ResourceCPU:     resource.MustParse("2000m"),
					corev1.ResourceMemory:  resource.MustParse("12Gi"),
					corev1.ResourceStorage: resource.MustParse("15Gi"),
				},
			},
		},
	}
	for _, node := range nodes {
		clientset.GetNodes().Create(context.TODO(), node, metav1.CreateOptions{})
	}

	namespaces := []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
		},
	}

	for _, namespace := range namespaces {
		clientset.GetNamespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
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
		{
			ObjectMeta: metav1.ObjectMeta{Name: "service-2", Namespace: "test"},
			Spec: corev1.ServiceSpec{
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "10.100.1.2",
				Selector:  map[string]string{"app": "test-app"},
			},
		},
	}
	for _, service := range services {
		namespace := service.GetNamespace()
		clientset.GetServices(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
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
		{
			ObjectMeta: metav1.ObjectMeta{Name: "deployment-2", Namespace: "test"},
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
		namespace := deployment.GetNamespace()
		clientset.GetDeployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	}

	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "secret-1", Namespace: "default"},
			Type:       corev1.SecretTypeOpaque,
			Data:       map[string][]byte{"password": []byte("supersecret")},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "secret-2", Namespace: "test"},
			Type:       corev1.SecretTypeOpaque,
			Data:       map[string][]byte{"password": []byte("supersecret")},
		},
	}
	for _, secret := range secrets {
		namespace := secret.GetNamespace()
		clientset.GetSecrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	}

	configMaps := []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "configmap-1", Namespace: "default"},
			Data:       map[string]string{"config": "valueaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "configmap-2", Namespace: "test"},
			Data:       map[string]string{"config": "value"},
		},
	}
	for _, configMap := range configMaps {
		namespace := configMap.GetNamespace()
		clientset.GetConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	}

	return clientset
}
