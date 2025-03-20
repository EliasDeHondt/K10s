/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/client"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sync"
	"time"
)

func GetSecretsHandler(ctx *gin.Context) {
	namespace, ok := ctx.GetQuery("namespace")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var secretList *PaginatedResponse[[]client.Secret]
	var err error

	if ok {
		secretList, err = GetSecrets(C, namespace, pageSize, pageToken)
	} else {
		secretList, err = GetSecrets(C, "", pageSize, pageToken)
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, secretList)
}

func GetSecrets(c client.IClient, namespace string, pageSize int, pageToken string) (*PaginatedResponse[[]client.Secret], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetSecrets(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]client.Secret]{
		Response:  transformSecrets(&list.Items),
		PageToken: list.Continue,
	}, nil
}

func transformSecrets(list *[]v1.Secret) []client.Secret {
	var secretList = make([]client.Secret, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, secret := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, secret v1.Secret) {
			defer wg.Done()
			secretList[i] = client.NewSecret(secret)
			<-semaphore
		}(i, secret)

	}
	wg.Wait()

	return secretList
}
