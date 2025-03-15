/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package kubernetes

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	av1 "k8s.io/api/apps/v1"
	cv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/metrics/pkg/apis/metrics"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
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
	GetNamespaces() corev1.NamespaceInterface
	GetNodes() corev1.NodeInterface
	GetPods(namespace string) corev1.PodInterface
	GetServices(namespace string) corev1.ServiceInterface
	GetEndpoints(namespace string) corev1.EndpointsInterface
	GetConfigMaps(namespace string) corev1.ConfigMapInterface
	GetSecrets(namespace string) corev1.SecretInterface
	GetDeployments(namespace string) appsv1.DeploymentInterface
	GetReplicaSets(namespace string) appsv1.ReplicaSetInterface
	WatchUsage()
	CreateNamespace(namespace *cv1.Namespace) (Namespace, error)
	CreateNode(node *cv1.Node) (Node, error)
	CreatePod(pod *cv1.Pod) (Pod, error)
	CreateDeployment(deployment *av1.Deployment) (Deployment, error)
	CreateService(service *cv1.Service) (Service, error)
	CreateConfigMap(configMap *cv1.ConfigMap) (ConfigMap, error)
	CreateSecret(secret *cv1.Secret) (Secret, error)
	GetFilteredPods(namespace string, nodeName string, pageSize int, continueToken string) (*[]cv1.Pod, string, error)
	GetFilteredServices(namespace string, nodeName string, pageSize int, continueToken string) (*[]cv1.Service, string, error)
	GetFilteredDeployments(namespace string, nodeName string, pageSize int, continueToken string) (*[]av1.Deployment, string, error)
	AddMetricsConnection(conn *websocket.Conn)
	GetTotalUsage() (*Metrics, error)
}

type FakeClient struct {
	Client        *fake.Clientset
	MetricsClient *FakeMetricsClient
	MetricsConns  []*websocket.Conn
}

type Client struct {
	Client        *kubernetes.Clientset
	MetricsClient *metricsv.Clientset
	MetricsConns  []*websocket.Conn
}

func (client *FakeClient) GetNamespaces() corev1.NamespaceInterface {
	return client.Client.CoreV1().Namespaces()
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

func (client *FakeClient) GetReplicaSets(namespace string) appsv1.ReplicaSetInterface {
	return client.Client.AppsV1().ReplicaSets(namespace)
}

func (client *FakeClient) CreateNamespace(namespace *cv1.Namespace) (Namespace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newNamespace, err := client.GetNamespaces().Create(ctx, namespace, metav1.CreateOptions{})

	if err != nil {
		return Namespace{}, err
	}

	return NewNamespace(*newNamespace), err
}

func (client *FakeClient) CreateNode(node *cv1.Node) (Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	node, err := client.GetNodes().Create(ctx, node, metav1.CreateOptions{})

	if err != nil {
		return Node{}, err
	}

	return NewNode(*node, client), err
}

func (client *FakeClient) CreatePod(pod *cv1.Pod) (Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if pod.Namespace == "" {
		pod.Namespace = "default"
	}

	pod, err := client.GetPods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})

	if err != nil {
		return Pod{}, err
	}

	return NewPod(*pod, client), err
}

func (client *FakeClient) CreateDeployment(deployment *av1.Deployment) (Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if deployment.Namespace == "" {
		deployment.Namespace = "default"
	}

	deployment, err := client.GetDeployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})

	if err != nil {
		return Deployment{}, err
	}

	return NewDeployment(*deployment), err
}

func (client *FakeClient) CreateService(service *cv1.Service) (Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if service.Namespace == "" {
		service.Namespace = "default"
	}

	_, err := client.GetServices(service.Namespace).Create(ctx, service, metav1.CreateOptions{})

	if err != nil {
		return Service{}, err
	}

	return NewService(*service), err
}

func (client *FakeClient) CreateConfigMap(configMap *cv1.ConfigMap) (ConfigMap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if configMap.Namespace == "" {
		configMap.Namespace = "default"
	}

	_, err := client.GetConfigMaps(configMap.Namespace).Create(ctx, configMap, metav1.CreateOptions{})

	if err != nil {
		return ConfigMap{}, err
	}

	return NewConfigMap(*configMap), err
}

func (client *FakeClient) CreateSecret(secret *cv1.Secret) (Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if secret.Namespace == "" {
		secret.Namespace = "default"
	}

	_, err := client.GetSecrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})

	if err != nil {
		return Secret{}, err
	}

	return NewSecret(*secret), err
}

func (client *Client) GetNamespaces() corev1.NamespaceInterface {
	return client.Client.CoreV1().Namespaces()
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

func (client *Client) GetReplicaSets(namespace string) appsv1.ReplicaSetInterface {
	return client.Client.AppsV1().ReplicaSets(namespace)
}

func (client *FakeClient) WatchUsage() {
	ctx := context.Background()

	nodes, err := client.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return
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

	calculatedMetrics := &Metrics{CpuUsage: cpuUsagePercent, MemUsage: memUsagePercent, DiskUsage: totalDiskUsage, DiskCapacity: totalDisk}

	for _, conn := range client.MetricsConns {
		err := conn.WriteJSON(calculatedMetrics)
		if err != nil {
			fmt.Println(err)
			CloseConn(conn)
		}
	}
}

func (client *FakeClient) AddMetricsConnection(conn *websocket.Conn) {
	client.MetricsConns = append(client.MetricsConns, conn)
}

func (client *Client) AddMetricsConnection(conn *websocket.Conn) {
	client.MetricsConns = append(client.MetricsConns, conn)
}

func (client *Client) WatchUsage() {
	ctx := context.Background()

	watcher, err := client.MetricsClient.MetricsV1beta1().NodeMetricses().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}
	defer func() {
		log.Println("Metrics watcher stopped")
		watcher.Stop()
	}()

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added, watch.Modified:
			calculatedMetrics, err := client.GetTotalUsage()
			if err != nil {
				log.Println(err)
				continue
			}

			for _, conn := range client.MetricsConns {
				err := conn.WriteJSON(calculatedMetrics)
				if err != nil {
					fmt.Println(err)
					CloseConn(conn)
				}
			}
		case watch.Deleted:
		case watch.Error:
			fmt.Printf("error watching node calculatedMetrics: %v", event.Object)
		}
	}
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
		totalDisk += node.Status.Capacity.StorageEphemeral().Value()

		nodeMetrics, err := client.MetricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, node.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		totalCpuUsage += nodeMetrics.Usage.Cpu().MilliValue()
		totalMemUsage += nodeMetrics.Usage.Memory().Value()
		totalDiskUsage += nodeMetrics.Usage.Storage().Value()
		totalDiskUsage += nodeMetrics.Usage.StorageEphemeral().Value()
	}

	cpuUsagePercent := float64(totalCpuUsage) / float64(totalCpu) * 100
	memUsagePercent := float64(totalMemUsage) / float64(totalMem) * 100

	return &Metrics{CpuUsage: cpuUsagePercent, MemUsage: memUsagePercent, DiskUsage: totalDiskUsage, DiskCapacity: totalDisk}, nil
}

func (client *Client) CreateNamespace(namespace *cv1.Namespace) (Namespace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newNamespace, err := client.GetNamespaces().Create(ctx, namespace, metav1.CreateOptions{})

	if err != nil {
		return Namespace{}, err
	}

	return NewNamespace(*newNamespace), err
}

func (client *Client) CreateNode(node *cv1.Node) (Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	node, err := client.GetNodes().Create(ctx, node, metav1.CreateOptions{})

	if err != nil {
		return Node{}, err
	}

	return NewNode(*node, client), err
}

func (client *Client) CreatePod(pod *cv1.Pod) (Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if pod.Namespace == "" {
		pod.Namespace = "default"
	}

	pod, err := client.GetPods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})

	if err != nil {
		return Pod{}, err
	}

	return NewPod(*pod, client), err
}

func (client *Client) CreateDeployment(deployment *av1.Deployment) (Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if deployment.Namespace == "" {
		deployment.Namespace = "default"
	}

	deployment, err := client.GetDeployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})

	if err != nil {
		return Deployment{}, err
	}

	return NewDeployment(*deployment), err
}

func (client *Client) CreateService(service *cv1.Service) (Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if service.Namespace == "" {
		service.Namespace = "default"
	}

	_, err := client.GetServices(service.Namespace).Create(ctx, service, metav1.CreateOptions{})

	if err != nil {
		return Service{}, err
	}

	return NewService(*service), err
}

func (client *Client) CreateConfigMap(configMap *cv1.ConfigMap) (ConfigMap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if configMap.Namespace == "" {
		configMap.Namespace = "default"
	}

	_, err := client.GetConfigMaps(configMap.Namespace).Create(ctx, configMap, metav1.CreateOptions{})

	if err != nil {
		return ConfigMap{}, err
	}

	return NewConfigMap(*configMap), err
}

func (client *Client) CreateSecret(secret *cv1.Secret) (Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if secret.Namespace == "" {
		secret.Namespace = "default"
	}

	_, err := client.GetSecrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})

	if err != nil {
		return Secret{}, err
	}

	return NewSecret(*secret), err
}

func (client *FakeClient) GetFilteredPods(namespace string, nodeName string, pageSize int, continueToken string) (*[]cv1.Pod, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filteredPods []cv1.Pod
	var newContinueToken string

	list, err := client.Client.CoreV1().Pods(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: continueToken,
	})
	if err != nil {
		return nil, "", err
	}

	for _, pod := range list.Items {
		if nodeName == "" || pod.Spec.NodeName == nodeName {
			filteredPods = append(filteredPods, pod)
		}
	}
	newContinueToken = list.Continue
	return &filteredPods, newContinueToken, nil
}

func (client *Client) GetFilteredPods(namespace string, nodeName string, pageSize int, continueToken string) (*[]cv1.Pod, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filteredPods []cv1.Pod
	var newContinueToken string

	fieldSelector := ""
	if nodeName != "" {
		fieldSelector = "spec.nodeName=" + nodeName
	}

	list, err := client.GetPods(namespace).List(ct, metav1.ListOptions{
		Limit:         int64(pageSize),
		Continue:      continueToken,
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, "", err
	}

	filteredPods = list.Items
	newContinueToken = list.Continue
	return &filteredPods, newContinueToken, nil
}

func (client *FakeClient) GetFilteredServices(namespace string, nodeName string, pageSize int, continueToken string) (*[]cv1.Service, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filteredServices []cv1.Service
	var newContinueToken string

	services, err := client.GetServices(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: continueToken,
	})
	if err != nil {
		return nil, "", err
	}

	if nodeName != "" {
		pods, podToken, err := client.GetFilteredPods(namespace, nodeName, 0, "")
		if err != nil {
			return nil, "", err
		}
		allPods := *pods
		for podToken != "" {
			extraPods, nextPodToken, err := client.GetFilteredPods(namespace, nodeName, 0, podToken)
			if err != nil {
				return nil, "", err
			}
			allPods = append(allPods, *extraPods...)

			podToken = nextPodToken
		}

		for _, pod := range allPods {
			for _, svc := range services.Items {
				if isPodMatchingService(pod, svc) {
					filteredServices = append(filteredServices, svc)
				}
			}
		}
	} else {
		filteredServices = services.Items
	}

	return &filteredServices, newContinueToken, nil
}

func (client *Client) GetFilteredServices(namespace string, nodeName string, pageSize int, continueToken string) (*[]cv1.Service, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filteredServices []cv1.Service
	var newContinueToken string

	services, err := client.GetServices(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: continueToken,
	})
	if err != nil {
		return nil, "", err
	}

	if nodeName != "" {
		pods, podToken, err := client.GetFilteredPods(namespace, nodeName, 0, "")
		if err != nil {
			return nil, "", err
		}
		allPods := *pods
		for podToken != "" {
			extraPods, nextPodToken, err := client.GetFilteredPods(namespace, nodeName, 0, podToken)
			if err != nil {
				return nil, "", err
			}
			allPods = append(allPods, *extraPods...)

			podToken = nextPodToken
		}

		for _, pod := range allPods {
			for _, svc := range services.Items {
				if isPodMatchingService(pod, svc) {
					filteredServices = append(filteredServices, svc)
				}
			}
		}
	} else {
		filteredServices = services.Items
	}

	return &filteredServices, newContinueToken, nil
}

func (client *FakeClient) GetFilteredDeployments(namespace string, nodeName string, pageSize int, continueToken string) (*[]av1.Deployment, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filteredDeployments []av1.Deployment
	var newContinueToken string

	if nodeName != "" {
		pods, podToken, err := client.GetFilteredPods(namespace, nodeName, 0, "")
		if err != nil {
			return nil, "", err
		}
		allPods := *pods
		for podToken != "" {
			extraPods, nextPodToken, err := client.GetFilteredPods(namespace, nodeName, 0, podToken)
			if err != nil {
				return nil, "", err
			}
			allPods = append(allPods, *extraPods...)

			podToken = nextPodToken
		}
		for _, pod := range allPods {
			for _, owner := range pod.OwnerReferences {
				if owner.Kind == "ReplicaSet" {
					rs, err := client.Client.AppsV1().ReplicaSets(namespace).Get(ct, owner.Name, metav1.GetOptions{})
					if err != nil {
						continue
					}

					for _, rsOwner := range rs.OwnerReferences {
						if rsOwner.Kind == "Deployment" {
							deployment, err := client.GetDeployments(namespace).Get(ct, rsOwner.Name, metav1.GetOptions{})
							if err == nil {
								filteredDeployments = append(filteredDeployments, *deployment)
							}
						}
					}
				}
			}
		}
	} else {
		deployments, err := client.GetDeployments(namespace).List(ct, metav1.ListOptions{
			Limit:    int64(pageSize),
			Continue: continueToken,
		})
		if err != nil {
			return nil, "", err
		}
		filteredDeployments = deployments.Items
	}

	return &filteredDeployments, newContinueToken, nil
}

func (client *Client) GetFilteredDeployments(namespace string, nodeName string, pageSize int, continueToken string) (*[]av1.Deployment, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filteredDeployments []av1.Deployment
	var newContinueToken string

	if nodeName != "" {
		pods, podToken, err := client.GetFilteredPods(namespace, nodeName, 0, "")
		if err != nil {
			return nil, "", err
		}
		allPods := *pods
		for podToken != "" {
			extraPods, nextPodToken, err := client.GetFilteredPods(namespace, nodeName, 0, podToken)
			if err != nil {
				return nil, "", err
			}
			allPods = append(allPods, *extraPods...)

			podToken = nextPodToken
		}
		for _, pod := range allPods {
			for _, owner := range pod.OwnerReferences {
				if owner.Kind == "ReplicaSet" {
					rs, err := client.Client.AppsV1().ReplicaSets(namespace).Get(ct, owner.Name, metav1.GetOptions{})
					if err != nil {
						continue
					}

					for _, rsOwner := range rs.OwnerReferences {
						if rsOwner.Kind == "Deployment" {
							deployment, err := client.GetDeployments(namespace).Get(ct, rsOwner.Name, metav1.GetOptions{})
							if err == nil {
								filteredDeployments = append(filteredDeployments, *deployment)
							}
						}
					}
				}
			}
		}
	} else {
		deployments, err := client.GetDeployments(namespace).List(ct, metav1.ListOptions{
			Limit:    int64(pageSize),
			Continue: continueToken,
		})
		if err != nil {
			return nil, "", err
		}
		filteredDeployments = deployments.Items
	}

	return &filteredDeployments, newContinueToken, nil
}
