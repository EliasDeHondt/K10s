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

func GetNodesHandler(ctx *gin.Context) {
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	nodeList, err := GetNodes(C, pageSize, pageToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, nodeList)
}

func GetNodes(c client.IClient, pageSize int, pageToken string) (*PaginatedResponse[[]client.Node], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetNodes().List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]client.Node]{
		Response:  transformNodes(&list.Items),
		PageToken: list.Continue,
	}, nil
}

func GetNodeNamesHandler(ctx *gin.Context) {
	names, err := GetNodeNames(C)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}

	ctx.JSON(http.StatusOK, names)
}

func GetNodeNames(c client.IClient) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	list, err := c.GetNodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return transformNodesToName(&list.Items), nil
}

func transformNodesToName(list *[]v1.Node) []string {
	result := make([]string, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, node := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, node v1.Node) {
			defer wg.Done()
			result[i] = node.Name
			<-semaphore
		}(i, node)
	}

	wg.Wait()
	return result
}

func transformNodes(list *[]v1.Node) []client.Node {
	var nodeList = make([]client.Node, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, node := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, node v1.Node) {
			defer wg.Done()
			nodeList[i] = client.NewNode(node, C)
			<-semaphore
		}(i, node)
	}

	wg.Wait()
	return nodeList
}
