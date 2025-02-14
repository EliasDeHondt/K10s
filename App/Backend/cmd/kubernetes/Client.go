package kubernetes

import (
	"context"
	cv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/metrics/pkg/apis/metrics"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"time"
)

type Metrics struct {
	CpuUsage     float64
	MemUsage     float64
	DiskUsage    int64
	DiskCapacity int64
}

type FakeMetricsClient struct {
	NodeMetrics map[string]*metrics.NodeMetrics
	PodMetrics  map[string]*metrics.PodMetrics
}

type IClient interface {
	GetNodes() corev1.NodeInterface
	GetPods(namespace string) corev1.PodInterface
	GetServices(namespace string) corev1.ServiceInterface
	GetEndpoints(namespace string) corev1.EndpointsInterface
	GetConfigMaps(namespace string) corev1.ConfigMapInterface
	GetSecrets(namespace string) corev1.SecretInterface
	GetDeployments(namespace string) appsv1.DeploymentInterface
	GetTotalUsage() (*Metrics, error)
	GetUsageForNode(nodeName string) (*Metrics, error)
	CreatePod(podInterface *cv1.Pod) error
}

type FakeClient struct {
	Client        *fake.Clientset
	MetricsClient *FakeMetricsClient
}

type Client struct {
	Client        *kubernetes.Clientset
	MetricsClient *metricsv.Clientset
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

func (client *FakeClient) GetConfigMaps(namespace string) corev1.ConfigMapInterface {
	return client.Client.CoreV1().ConfigMaps(namespace)
}

func (client *FakeClient) GetSecrets(namespace string) corev1.SecretInterface {
	return client.Client.CoreV1().Secrets(namespace)
}

func (client *FakeClient) GetDeployments(namespace string) appsv1.DeploymentInterface {
	return client.Client.AppsV1().Deployments(namespace)
}

func (client *FakeClient) CreatePod(pod *cv1.Pod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.GetPods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	return err
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

func (client *Client) GetConfigMaps(namespace string) corev1.ConfigMapInterface {
	return client.Client.CoreV1().ConfigMaps(namespace)
}

func (client *Client) GetSecrets(namespace string) corev1.SecretInterface {
	return client.Client.CoreV1().Secrets(namespace)
}

func (client *Client) GetDeployments(namespace string) appsv1.DeploymentInterface {
	return client.Client.AppsV1().Deployments(namespace)
}

func (client *FakeClient) GetTotalUsage() (*Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodes, err := client.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var totalCpu int64
	var totalMem int64
	var totalDisk int64
	var totalCpuUsage int64
	var totalMemUsage int64
	var totalDiskUsage int64

	for _, node := range nodes.Items {
		totalCpu += node.Status.Capacity.Cpu().MilliValue()
		totalMem += node.Status.Capacity.Memory().Value()
		totalDisk += node.Status.Capacity.Storage().Value()

		nodeMetrics := client.MetricsClient.NodeMetrics[node.Name]

		totalCpuUsage += nodeMetrics.Usage.Cpu().MilliValue()
		totalMemUsage += nodeMetrics.Usage.Memory().Value()
		totalDiskUsage += nodeMetrics.Usage.Storage().Value()
	}

	cpuUsagePercent := float64(totalCpuUsage) / float64(totalCpu) * 100
	memUsagePercent := float64(totalMemUsage) / float64(totalMem) * 100

	return &Metrics{CpuUsage: cpuUsagePercent, MemUsage: memUsagePercent, DiskUsage: totalDiskUsage, DiskCapacity: totalDisk}, nil
}

func (client *FakeClient) GetUsageForNode(nodeName string) (*Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	node, err := client.Client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	totalCpu := node.Status.Capacity.Cpu().MilliValue()
	totalMem := node.Status.Capacity.Memory().Value()
	totalDisk := node.Status.Capacity.Storage().Value()

	nodeMetrics := client.MetricsClient.NodeMetrics[nodeName]

	usedCpu := nodeMetrics.Usage.Cpu().MilliValue()
	usedMem := nodeMetrics.Usage.Memory().Value()
	usedDisk := nodeMetrics.Usage.Storage().Value()

	cpuUsagePercent := float64(usedCpu) / float64(totalCpu) * 100
	memoryUsagePercent := float64(usedMem) / float64(totalMem) * 100

	return &Metrics{CpuUsage: cpuUsagePercent, MemUsage: memoryUsagePercent, DiskUsage: usedDisk, DiskCapacity: totalDisk}, nil
}

func (client *Client) GetTotalUsage() (*Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodes, err := client.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var totalCpu int64
	var totalMem int64
	var totalDisk int64
	var totalCpuUsage int64
	var totalMemUsage int64
	var totalDiskUsage int64

	for _, node := range nodes.Items {
		totalCpu += node.Status.Capacity.Cpu().MilliValue()
		totalMem += node.Status.Capacity.Memory().Value()
		totalDisk += node.Status.Capacity.Storage().Value()

		nodeMetrics, err := client.MetricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, node.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		totalCpuUsage += nodeMetrics.Usage.Cpu().MilliValue()
		totalMemUsage += nodeMetrics.Usage.Memory().Value()
		totalDiskUsage += nodeMetrics.Usage.Storage().Value()
	}

	cpuUsagePercent := float64(totalCpuUsage) / float64(totalCpu) * 100
	memUsagePercent := float64(totalMemUsage) / float64(totalMem) * 100

	return &Metrics{CpuUsage: cpuUsagePercent, MemUsage: memUsagePercent, DiskUsage: totalDiskUsage, DiskCapacity: totalDisk}, nil
}

func (client *Client) GetUsageForNode(nodeName string) (*Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	node, err := client.Client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	totalCpu := node.Status.Capacity.Cpu().MilliValue()
	totalMem := node.Status.Capacity.Memory().Value()
	totalDisk := node.Status.Capacity.Storage().Value()

	nodeMetrics, err := client.MetricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	usedCpu := nodeMetrics.Usage.Cpu().MilliValue()
	usedMem := nodeMetrics.Usage.Memory().Value()
	usedDisk := nodeMetrics.Usage.Storage().Value()

	cpuUsagePercent := float64(usedCpu) / float64(totalCpu) * 100
	memoryUsagePercent := float64(usedMem) / float64(totalMem) * 100

	return &Metrics{CpuUsage: cpuUsagePercent, MemUsage: memoryUsagePercent, DiskUsage: usedDisk, DiskCapacity: totalDisk}, nil
}
