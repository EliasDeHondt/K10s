package kubernetes

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"time"
)

type Service struct {
	Namespace  string
	Name       string
	Type       string
	ClusterIp  string
	ExternalIp []string
	Ports      []int
	Age        string
}

func NewService(service v1.Service) Service {
	return Service{
		Namespace:  service.Namespace,
		Name:       service.Name,
		Type:       string(service.Spec.Type),
		ClusterIp:  service.Spec.ClusterIP,
		ExternalIp: service.Spec.ExternalIPs,
		Ports:      getPorts(&service),
		Age:        calculateServiceAge(&service),
	}
}

func getPorts(service *v1.Service) []int {
	var portList = make([]int, 0)
	ports := service.Spec.Ports

	for _, port := range ports {
		portList = append(portList, int(port.Port))
	}
	return portList
}

func calculateServiceAge(service *v1.Service) string {
	age := time.Since(service.CreationTimestamp.Time)

	return fmt.Sprintf("%d:%d:%d", int(age.Hours()/24), int(age.Hours())%24, int(age.Minutes())%60)
}
