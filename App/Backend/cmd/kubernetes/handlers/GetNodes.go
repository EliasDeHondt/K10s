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

func GetNodesHandler(ctx *gin.Context) {
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	nodeList, err := GetNodes(c, pageSize, pageToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, nodeList)
}

func GetNodes(c kubernetes.IClient, pageSize int, pageToken string) (*PaginatedResponse[[]kubernetes.Node], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetNodes().List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]kubernetes.Node]{
		Response:  transformNodes(&list.Items),
		PageToken: list.Continue,
	}, nil
}

func transformNodes(list *[]v1.Node) []kubernetes.Node {
	var nodeList = make([]kubernetes.Node, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, node := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, node v1.Node) {
			defer wg.Done()
			nodeList[i] = kubernetes.NewNode(node, c)
			<-semaphore
		}(i, node)
	}

	wg.Wait()
	return nodeList
}
