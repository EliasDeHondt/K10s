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

func GetServicesHandler(ctx *gin.Context) {
	namespace, ok := ctx.GetQuery("namespace")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var serviceList *PaginatedResponse[[]kubernetes.Service]
	var err error

	if ok {
		serviceList, err = GetServices(c, namespace, pageSize, pageToken)
	} else {
		serviceList, err = GetServices(c, "", pageSize, pageToken)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, serviceList)
}

func GetServices(c kubernetes.IClient, namespace string, pageSize int, pageToken string) (*PaginatedResponse[[]kubernetes.Service], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetServices(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]kubernetes.Service]{
		Response:  transformServices(&list.Items),
		PageToken: list.Continue,
	}, nil
}

func transformServices(list *[]v1.Service) []kubernetes.Service {
	var serviceList = make([]kubernetes.Service, len(*list))
	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, service := range *list {
		wg.Add(1)

		go func(i int, service v1.Service) {
			defer wg.Done()
			serviceList[i] = kubernetes.NewService(service)
			semaphore <- struct{}{}
		}(i, service)
	}
	wg.Wait()

	return serviceList
}
