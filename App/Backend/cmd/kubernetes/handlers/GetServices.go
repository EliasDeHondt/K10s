package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetServicesHandler(ctx *gin.Context) {
	podList, err := getServices(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, podList)
}

func getServices(c kubernetes.IClient) (*[]kubernetes.Service, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetServices("").List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var serviceList = make([]kubernetes.Service, len(list.Items))

	for _, service := range list.Items {
		serviceList = append(serviceList, kubernetes.NewService(service))
	}
	return &serviceList, nil
}
