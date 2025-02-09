package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type IClient interface {
	GetNodes() corev1.NodeInterface
	GetPods(namespace string) corev1.PodInterface
	GetServices(namespace string) corev1.ServiceInterface
	GetEndpoints(namespace string) corev1.EndpointsInterface
}

type FakeClient struct {
	Client *fake.Clientset
}

type Client struct {
	Client *kubernetes.Clientset
}

func (client *FakeClient) GetNodes() corev1.NodeInterface {
	return client.Client.CoreV1().Nodes()
}

func (client *FakeClient) GetPods(namespace string) corev1.PodInterface {
	return client.Client.CoreV1().Pods(namespace)
}

func (client *FakeClient) GetServices(namespace string) corev1.ServiceInterface {
	return client.Client.CoreV1().Services(namespace)
}

func (client *FakeClient) GetEndpoints(namespace string) corev1.EndpointsInterface {
	return client.Client.CoreV1().Endpoints(namespace)
}

func (client *Client) GetNodes() corev1.NodeInterface {
	return client.Client.CoreV1().Nodes()
}

func (client *Client) GetPods(namespace string) corev1.PodInterface {
	return client.Client.CoreV1().Pods(namespace)
}

func (client *Client) GetServices(namespace string) corev1.ServiceInterface {
	return client.Client.CoreV1().Services(namespace)
}

func (client *Client) GetEndpoints(namespace string) corev1.EndpointsInterface {
	return client.Client.CoreV1().Endpoints(namespace)
}
