package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetPodsHandler(ctx *gin.Context) {
	podList, err := GetPods(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, podList)
}

func GetPods(c kubernetes.IClient) (*[]kubernetes.Pod, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetPods("").List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podList = make([]kubernetes.Pod, len(list.Items))

	for i, pod := range list.Items {
		podList[i] = kubernetes.NewPod(pod, c)
	}
	return &podList, nil
}
