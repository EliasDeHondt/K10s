/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/client"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	"net/http"
	"sync"
)

func GetDeploymentsHandler(ctx *gin.Context) {
	namespace, _ := ctx.GetQuery("namespace")
	nodeName, _ := ctx.GetQuery("node")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var deploymentList *PaginatedResponse[[]client.Deployment]
	var err error

	deploymentList, err = GetDeployments(C, namespace, nodeName, pageSize, pageToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, deploymentList)
}

func GetDeployments(c client.IClient, namespace string, nodeName string, pageSize int, pageToken string) (*PaginatedResponse[[]client.Deployment], error) {
	list, token, err := c.GetFilteredDeployments(namespace, nodeName, pageSize, pageToken)
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]client.Deployment]{
		Response:  transformDeployments(list),
		PageToken: token,
	}, nil
}

func transformDeployments(list *[]v1.Deployment) []client.Deployment {
	var deploymentList = make([]client.Deployment, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, deployment := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, deployment v1.Deployment) {
			defer wg.Done()
			deploymentList[i] = client.NewDeployment(deployment)
			<-semaphore
		}(i, deployment)
	}

	wg.Wait()
	return deploymentList
}
