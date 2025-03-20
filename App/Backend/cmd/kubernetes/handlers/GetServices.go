/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/client"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"sync"
)

func GetServicesHandler(ctx *gin.Context) {
	namespace, _ := ctx.GetQuery("namespace")
	nodeName, _ := ctx.GetQuery("node")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	serviceList, err := GetServices(C, namespace, nodeName, pageSize, pageToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, serviceList)
}

func GetServices(c client.IClient, namespace string, nodeName string, pageSize int, pageToken string) (*PaginatedResponse[[]client.Service], error) {
	list, token, err := c.GetFilteredServices(namespace, nodeName, pageSize, pageToken)
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]client.Service]{
		Response:  transformServices(list),
		PageToken: token,
	}, nil
}

func transformServices(list *[]v1.Service) []client.Service {
	var serviceList = make([]client.Service, len(*list))
	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, service := range *list {
		wg.Add(1)

		go func(i int, service v1.Service) {
			defer wg.Done()
			serviceList[i] = client.NewService(service)
			semaphore <- struct{}{}
		}(i, service)
	}
	wg.Wait()

	return serviceList
}
