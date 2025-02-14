package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sync"
	"time"
)

func GetDeploymentsHandler(ctx *gin.Context) {
	namespace, ok := ctx.GetQuery("namespace")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var deploymentList *PaginatedResponse[[]kubernetes.Deployment]
	var err error

	if ok {
		deploymentList, err = GetDeployments(c, namespace, pageSize, pageToken)
	} else {
		deploymentList, err = GetDeployments(c, "", pageSize, pageToken)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, deploymentList)
}

func GetDeployments(c *kubernetes.IClient, namespace string, pageSize int, pageToken string) (*PaginatedResponse[[]kubernetes.Deployment], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := (*c).GetDeployments(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}

	var deploymentList = make([]kubernetes.Deployment, len(list.Items))

	for i, deployment := range list.Items {
		deploymentList[i] = kubernetes.NewDeployment(deployment)
	}
	return &PaginatedResponse[[]kubernetes.Deployment]{
		Response:  transformDeployments(&list.Items),
		PageToken: list.Continue,
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

		go func(i int, node v1.Deployment) {
			defer wg.Done()
			deploymentList[i] = kubernetes.NewDeployment(deployment)
			<-semaphore
		}(i, deployment)
	}

	wg.Wait()
	return deploymentList
}
