/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sync"
	"time"
)

func GetNamespacesHandler(ctx *gin.Context) {

	namespaces, err := GetNamespaces(c)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}

	ctx.JSON(http.StatusOK, namespaces)
}

func GetNamespaces(c kubernetes.IClient) ([]kubernetes.Namespace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	namespaces, err := c.GetNamespaces().List(ctx, metav1.ListOptions{
		Limit: 0,
	})

	if err != nil {
		return nil, err
	}

	return transformNamespaces(&namespaces.Items), nil
}

func transformNamespaces(list *[]v1.Namespace) []kubernetes.Namespace {
	var namespaceList = make([]kubernetes.Namespace, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, namespace := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, namespace v1.Namespace) {
			defer wg.Done()
			namespaceList[i] = kubernetes.NewNamespace(namespace)
			<-semaphore
		}(i, namespace)
	}

	wg.Wait()
	return namespaceList
}
