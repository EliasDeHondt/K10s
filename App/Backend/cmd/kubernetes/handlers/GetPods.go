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

func GetPodsHandler(ctx *gin.Context) {
	namespace, _ := ctx.GetQuery("namespace")
	nodeName, _ := ctx.GetQuery("node")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var podList *PaginatedResponse[[]client.Pod]
	var err error

	podList, err = GetPods(C, namespace, nodeName, pageSize, pageToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, podList)
}

func GetPods(c client.IClient, namespace string, nodeName string, pageSize int, continueToken string) (*PaginatedResponse[[]client.Pod], error) {

	filteredPods, newContinueToken, err := c.GetFilteredPods(namespace, nodeName, pageSize, continueToken)
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]client.Pod]{
		Response:  transformPods(filteredPods),
		PageToken: newContinueToken,
	}, nil
}

func transformPods(list *[]v1.Pod) []client.Pod {
	var podList = make([]client.Pod, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, pod := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, pod v1.Pod) {
			defer wg.Done()
			podList[i] = client.NewPod(pod, C)
			<-semaphore
		}(i, pod)
	}

	wg.Wait()
	return podList
}
