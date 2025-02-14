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
	namespace, ok := ctx.GetQuery("namespace")

	var serviceList *[]kubernetes.Service
	var err error

	if ok {
		serviceList, err = GetServices(c, namespace)
	} else {
		serviceList, err = GetServices(c, "")
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, serviceList)
}

func GetServices(c kubernetes.IClient, namespace string) (*[]kubernetes.Service, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetServices(namespace).List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var serviceList = make([]kubernetes.Service, len(list.Items))

	for i, service := range list.Items {
		serviceList[i] = kubernetes.NewService(service)
	}
	return &serviceList, nil
}
