package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetSecretsHandler(ctx *gin.Context) {
	nodeList, err := getSecrets(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, nodeList)
}

func getSecrets(c kubernetes.IClient) (*[]kubernetes.Secret, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetSecrets("").List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var secretList = make([]kubernetes.Secret, len(list.Items))

	for i, secret := range list.Items {
		secretList[i] = kubernetes.NewSecret(secret)
	}
	return &secretList, nil
}
