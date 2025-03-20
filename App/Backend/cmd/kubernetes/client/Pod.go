/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package client

import (
	"context"
	"fmt"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
)

type Pod struct {
	Namespace     string
	Name          string
	ServicesReady string
	Restarts      int
	Status        string
	IP            string
	Node          string
	Age           string
}

func NewPod(pod v1.Pod, clientset IClient) Pod {

	runningServices, totalServices := getServiceHealthForPod(pod, clientset)

	readyString := fmt.Sprintf("%d/%d", runningServices, totalServices)

	return Pod{
		Namespace:     pod.Namespace,
		Name:          pod.Name,
		ServicesReady: readyString,
		Restarts:      getTotalContainerRestarts(pod),
		Status:        isPodOnline(&pod),
		IP:            pod.Status.PodIP,
		Node:          pod.Spec.NodeName,
		Age:           calculatePodAge(pod),
	}
}

func isPodOnline(pod *v1.Pod) string {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodReady {
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

func getServiceHealthForPod(pod v1.Pod, clientset IClient) (int, int) {
	services := getServicesForPod(pod, clientset)
	totalServices := len(services)
	runningServices := 0

	for _, serviceName := range services {
		if isServiceRunning(pod.Namespace, serviceName, clientset) {
			runningServices++
		}
	}

	return runningServices, totalServices
}

func getServicesForPod(pod v1.Pod, clientset IClient) []string {
	var matchingServices []string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	services, err := clientset.GetServices(pod.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing services: %v", err)
		return nil
	}

	for _, service := range services.Items {
		if util.IsPodMatchingService(pod, service) {
			matchingServices = append(matchingServices, service.Name)
		}
	}

	return matchingServices
}

func isServiceRunning(namespace, serviceName string, clientset IClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	endpoints, err := clientset.GetEndpoints(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting endpoints for service %s: %v", serviceName, err)
		return false
	}

	for _, subset := range endpoints.Subsets {
		if len(subset.Addresses) > 0 {
			return true
		}
	}

	return false
}

func getTotalContainerRestarts(pod v1.Pod) int {
	totalRestarts := 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		totalRestarts += int(containerStatus.RestartCount)
	}
	return totalRestarts
}

func calculatePodAge(pod v1.Pod) string {
	age := time.Since(pod.CreationTimestamp.Time)

	return fmt.Sprintf("%d:%d:%d", int(age.Hours()/24), int(age.Hours())%24, int(age.Minutes())%60)
}
