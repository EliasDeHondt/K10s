/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package client

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"time"
)

type ConfigMap struct {
	Namespace string
	Name      string
	Data      map[string]string
	Age       string
}

func NewConfigMap(confMap v1.ConfigMap) ConfigMap {
	return ConfigMap{
		Namespace: confMap.Namespace,
		Name:      confMap.Name,
		Data:      confMap.Data,
		Age:       getConfigMapAge(&confMap),
	}
}

func getConfigMapAge(confMap *v1.ConfigMap) string {
	age := time.Since(confMap.CreationTimestamp.Time)

	return fmt.Sprintf("%d:%d:%d", int(age.Hours()/24), int(age.Hours())%24, int(age.Minutes())%60)
}
