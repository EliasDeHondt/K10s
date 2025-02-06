package handlers

import (
	"context"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"net/http"
	"time"
)

var c = kubernetes.TestFakeClient()

func GetNodesHandler(ctx *gin.Context) {
	nodeList, err := getNodes(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nodeList)
}

func getNodes(c *fake.Clientset) (*[]kubernetes.Node, error) {
	ct, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := c.CoreV1().Nodes().List(ct, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeList []kubernetes.Node

	for _, node := range list.Items {
		nodeList = append(nodeList, kubernetes.NewNode(node))
	}
	return &nodeList, nil
}
