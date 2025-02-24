/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"sync"
)

func GetServicesHandler(ctx *gin.Context) {
	namespace, _ := ctx.GetQuery("namespace")
	nodeName, _ := ctx.GetQuery("node")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	serviceList, err := GetServices(c, namespace, nodeName, pageSize, pageToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, serviceList)
}

func GetServices(c kubernetes.IClient, namespace string, nodeName string, pageSize int, pageToken string) (*PaginatedResponse[[]kubernetes.Service], error) {
	/*ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetServices(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}*/

	list, token, err := c.GetFilteredServices(namespace, nodeName, pageSize, pageToken)
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]kubernetes.Service]{
		Response:  transformServices(list),
		PageToken: token,
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
