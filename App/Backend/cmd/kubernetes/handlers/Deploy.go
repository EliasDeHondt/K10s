package handlers

import (
	"fmt"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
)

func DeploymentHandler(ctx *gin.Context) {
	contentType := ctx.Request.Header.Get("Content-Type")
	if contentType != "application/x-yaml" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Content-Type"})
		return
	}

	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = Deployment(c, data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func Deployment(c *kubernetes.IClient, data []byte) error {
	decoder := scheme.Codecs.UniversalDeserializer()

	obj, gvk, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return err
	}

	switch gvk.Kind {
	case "Pod":
		pod := obj.(*corev1.Pod)
		err = (*c).CreatePod(pod)
	default:
		return fmt.Errorf("unknown kind: %s", gvk.Kind)
	}

	return err
}
