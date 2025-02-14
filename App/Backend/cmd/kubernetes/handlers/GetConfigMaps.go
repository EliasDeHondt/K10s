package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetConfigMapsHandler(ctx *gin.Context) {
	namespace, ok := ctx.GetQuery("namespace")

	var configMapList *[]kubernetes.ConfigMap
	var err error

	if ok {
		configMapList, err = GetConfigMaps(c, namespace)
	} else {
		configMapList, err = GetConfigMaps(c, "")
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, configMapList)
}

func GetConfigMaps(c *kubernetes.IClient, namespace string) (*[]kubernetes.ConfigMap, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := (*c).GetConfigMaps(namespace).List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var configMapList = make([]kubernetes.ConfigMap, len(list.Items))

	for i, configMap := range list.Items {
		configMapList[i] = kubernetes.NewConfigMap(configMap)
	}
	return &configMapList, nil
}
