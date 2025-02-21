/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package kubernetes

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"time"
)

type Deployment struct {
	Namespace string
	Name      string
	Ready     string
	Updated   bool
	Available bool
	Age       string
}

func NewDeployment(deployment v1.Deployment) Deployment {
	isUpdated := deployment.Status.UpdatedReplicas == deployment.Status.Replicas
	isAvailable := deployment.Status.AvailableReplicas == deployment.Status.Replicas

	return Deployment{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
		Ready:     getReadyPods(&deployment),
		Updated:   isUpdated,
		Available: isAvailable,
		Age:       getDeploymentAge(&deployment),
	}
}

func getReadyPods(deployment *v1.Deployment) string {
	totalPods := deployment.Status.Replicas
	readyPods := deployment.Status.ReadyReplicas

	return fmt.Sprintf("%d/%d", readyPods, totalPods)
}

func getDeploymentAge(deployment *v1.Deployment) string {
	age := time.Since(deployment.ObjectMeta.CreationTimestamp.Time)

	return fmt.Sprintf("%d:%d:%d", int(age.Hours()/24), int(age.Hours())%24, int(age.Minutes())%60)
}