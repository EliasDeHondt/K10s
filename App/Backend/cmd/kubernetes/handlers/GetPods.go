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
	namespace, _ := ctx.GetQuery("namespace")
	nodeName, _ := ctx.GetQuery("node")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var podList *PaginatedResponse[[]kubernetes.Pod]
	var err error

	podList, err = GetPods(c, namespace, nodeName, pageSize, pageToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, podList)
}

func GetPods(c kubernetes.IClient, namespace string, nodeName string, pageSize int, continueToken string) (*PaginatedResponse[[]kubernetes.Pod], error) {
	var filteredPods *[]v1.Pod
	var newContinueToken string
	var err error

	if fakeClient, ok := c.(*kubernetes.FakeClient); ok {
		filteredPods, newContinueToken, err = getFilteredPodsForFakeClient(fakeClient, namespace, nodeName)
		if err != nil {
			return nil, err
		}

	} else {
		filteredPods, newContinueToken, err = getFilteredPodsForRealClient(c, namespace, nodeName, pageSize, continueToken)
	}

	return &PaginatedResponse[[]kubernetes.Pod]{
		Response:  transformPods(filteredPods),
		PageToken: newContinueToken,
	}, nil
}

func getFilteredPodsForRealClient(client kubernetes.IClient, namespace string, nodeName string, pageSize int, continueToken string) (*[]v1.Pod, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filteredPods []v1.Pod
	var newContinueToken string

	fieldSelector := ""
	if nodeName != "" {
		fieldSelector = "spec.nodeName=" + nodeName
	}

	list, err := client.GetPods(namespace).List(ct, metav1.ListOptions{
		Limit:         int64(pageSize),
		Continue:      continueToken,
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, "", err
	}

	filteredPods = list.Items
	newContinueToken = list.Continue
	return &filteredPods, newContinueToken, nil
}

func getFilteredPodsForFakeClient(fakeClient *kubernetes.FakeClient, namespace string, nodeName string) (*[]v1.Pod, string, error) {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filteredPods []v1.Pod
	var newContinueToken string

	list, err := fakeClient.Client.CoreV1().Pods(namespace).List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, "", err
	}

	for _, pod := range list.Items {
		if nodeName == "" || pod.Spec.NodeName == nodeName {
			filteredPods = append(filteredPods, pod)
		}
	}
	newContinueToken = list.Continue
	return &filteredPods, newContinueToken, nil
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
