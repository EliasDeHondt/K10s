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
	configMapList, err := GetConfigMaps(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, configMapList)
}

func GetConfigMaps(c *kubernetes.IClient) (*[]kubernetes.ConfigMap, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := (*c).GetConfigMaps("").List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var configMapList = make([]kubernetes.ConfigMap, len(list.Items))

	for i, configMap := range list.Items {
		configMapList[i] = kubernetes.NewConfigMap(configMap)
	}
	return &configMapList, nil
}
