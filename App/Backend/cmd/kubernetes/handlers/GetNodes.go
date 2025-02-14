package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

func GetNodesHandler(ctx *gin.Context) {
	nodeList, err := GetNodes(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error has occurred or the request has been timed out."})
		return
	}
	ctx.JSON(http.StatusOK, nodeList)
}

func GetNodes(c kubernetes.IClient) (*[]kubernetes.Node, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.GetNodes().List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeList = make([]kubernetes.Node, len(list.Items))

	for i, node := range list.Items {
		nodeList[i] = kubernetes.NewNode(node, c)
	}
	return &nodeList, nil
}
