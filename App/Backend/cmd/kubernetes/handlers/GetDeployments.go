/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	"net/http"
	"sync"
)

func GetDeploymentsHandler(ctx *gin.Context) {
	namespace, _ := ctx.GetQuery("namespace")
	nodeName, _ := ctx.GetQuery("node")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var deploymentList *PaginatedResponse[[]kubernetes.Deployment]
	var err error

	deploymentList, err = GetDeployments(C, namespace, nodeName, pageSize, pageToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, deploymentList)
}

func GetDeployments(c kubernetes.IClient, namespace string, nodeName string, pageSize int, pageToken string) (*PaginatedResponse[[]kubernetes.Deployment], error) {
	list, token, err := c.GetFilteredDeployments(namespace, nodeName, pageSize, pageToken)
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]kubernetes.Deployment]{
		Response:  transformDeployments(list),
		PageToken: token,
	}, nil
}

func transformDeployments(list *[]v1.Deployment) []kubernetes.Deployment {
	var deploymentList = make([]kubernetes.Deployment, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, deployment := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, deployment v1.Deployment) {
			defer wg.Done()
			deploymentList[i] = kubernetes.NewDeployment(deployment)
			<-semaphore
		}(i, deployment)
	}

	wg.Wait()
	return deploymentList
}
