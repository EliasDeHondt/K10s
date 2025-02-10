package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetDeploymentsHandler(ctx *gin.Context) {
	deploymentList, err := GetDeployments(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, deploymentList)
}

func GetDeployments(c kubernetes.IClient) (*[]kubernetes.Deployment, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetDeployments("").List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentList = make([]kubernetes.Deployment, len(list.Items))

	for i, deployment := range list.Items {
		deploymentList[i] = kubernetes.NewDeployment(deployment)
	}
	return &deploymentList, nil
}
