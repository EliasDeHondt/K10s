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

func GetConfigMapsHandler(ctx *gin.Context) {
	namespace, ok := ctx.GetQuery("namespace")
	pageSize, pageToken := GetPageSizeAndPageToken(ctx)

	var configMapList *PaginatedResponse[[]client.ConfigMap]
	var err error

	if ok {
		configMapList, err = GetConfigMaps(C, namespace, pageSize, pageToken)
	} else {
		configMapList, err = GetConfigMaps(C, "", pageSize, pageToken)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, configMapList)
}

func GetConfigMaps(c client.IClient, namespace string, pageSize int, pageToken string) (*PaginatedResponse[[]client.ConfigMap], error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetConfigMaps(namespace).List(ct, metav1.ListOptions{
		Limit:    int64(pageSize),
		Continue: pageToken,
	})
	if err != nil {
		return nil, err
	}

	return &PaginatedResponse[[]client.ConfigMap]{
		Response:  transformConfigMaps(&list.Items),
		PageToken: list.Continue,
	}, nil
}

func transformConfigMaps(list *[]v1.ConfigMap) []client.ConfigMap {
	var configList = make([]client.ConfigMap, len(*list))

	var wg sync.WaitGroup
	concurrency := 20
	semaphore := make(chan struct{}, concurrency)

	for i, config := range *list {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, config v1.ConfigMap) {
			defer wg.Done()
			configList[i] = client.NewConfigMap(config)
			<-semaphore
		}(i, config)
	}

	wg.Wait()
	return configList
}
