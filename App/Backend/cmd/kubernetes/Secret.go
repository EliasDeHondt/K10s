/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package kubernetes

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"time"
)

type Secret struct {
	Namespace string
	Name      string
	Type      string
	Data      map[string][]byte
	Age       string
}

func NewSecret(secret v1.Secret) Secret {
	return Secret{
		Namespace: secret.Namespace,
		Name:      secret.Name,
		Type:      string(secret.Type),
		Data:      secret.Data,
		Age:       calculateSecretAge(&secret),
	}
}

func calculateSecretAge(secret *v1.Secret) string {
	age := time.Since(secret.CreationTimestamp.Time)

	return fmt.Sprintf("%d:%d:%d", int(age.Hours()/24), int(age.Hours())%24, int(age.Minutes())%60)
}