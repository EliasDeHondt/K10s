/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package client

import (
	"github.com/gorilla/websocket"
	av1 "k8s.io/api/apps/v1"
	cv1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Metrics struct {
	CpuUsage     float64
	MemUsage     float64
	DiskUsage    int64
	DiskCapacity int64
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
