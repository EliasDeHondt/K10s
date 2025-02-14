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

func GetConfigMapsHandler(ctx *gin.Context) {
	namespace, ok := ctx.GetQuery("namespace")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var configMapList *PaginatedResponse[[]kubernetes.ConfigMap]
	var err error

	if ok {
		configMapList, err = GetConfigMaps(c, namespace, pageSize, pageToken)
	} else {
		configMapList, err = GetConfigMaps(c, "", pageSize, pageToken)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, configMapList)
}

func GetConfigMaps(c *kubernetes.IClient, namespace string, pageSize int, pageToken string) (*PaginatedResponse[[]kubernetes.ConfigMap], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := (*c).GetConfigMaps(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}

	var configMapList = make([]kubernetes.ConfigMap, len(list.Items))

	for i, configMap := range list.Items {
		configMapList[i] = kubernetes.NewConfigMap(configMap)
	}
	return &PaginatedResponse[[]kubernetes.ConfigMap]{
		Response:  transformConfigMaps(&list.Items),
		PageToken: list.Continue,
	}, nil
}

func transformConfigMaps(list *[]v1.ConfigMap) []kubernetes.ConfigMap {
	var configList = make([]kubernetes.ConfigMap, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, config := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, node v1.ConfigMap) {
			defer wg.Done()
			configList[i] = kubernetes.NewConfigMap(config)
			<-semaphore
		}(i, config)
	}

	wg.Wait()
	return configList
}
