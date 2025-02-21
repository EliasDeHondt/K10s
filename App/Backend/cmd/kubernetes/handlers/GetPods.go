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

func GetPodsHandler(ctx *gin.Context) {
	namespace, ok := ctx.GetQuery("namespace")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var podList *PaginatedResponse[[]kubernetes.Pod]
	var err error

	if ok {
		podList, err = GetPods(c, namespace, pageSize, pageToken)
	} else {
		podList, err = GetPods(c, "", pageSize, pageToken)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, podList)
}

func GetPods(c kubernetes.IClient, namespace string, pageSize int, continueToken string) (*PaginatedResponse[[]kubernetes.Pod], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetPods(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: continueToken,
	})
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]kubernetes.Pod]{
		Response:  transformPods(&list.Items),
		PageToken: list.Continue,
	}, nil
}

func transformPods(list *[]v1.Pod) []kubernetes.Pod {
	var podList = make([]kubernetes.Pod, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, pod := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, pod v1.Pod) {
			defer wg.Done()
			podList[i] = kubernetes.NewPod(pod, c)
			<-semaphore
		}(i, pod)
	}

	wg.Wait()
	return podList
}