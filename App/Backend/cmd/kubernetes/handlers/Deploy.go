package handlers

import (
	"fmt"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
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
	case "Node":
		node := obj.(*corev1.Node)
		err = (*c).CreateNode(node)
	case "Pod":
		pod := obj.(*corev1.Pod)
		err = (*c).CreatePod(pod)
	case "Deployment":
		deployment := obj.(*appsv1.Deployment)
		err = (*c).CreateDeployment(deployment)
	case "Service":
		service := obj.(*corev1.Service)
		err = (*c).CreateService(service)
	case "ConfigMap":
		configMap := obj.(*corev1.ConfigMap)
		err = (*c).CreateConfigMap(configMap)
	case "Secret":
		secret := obj.(*corev1.Secret)
		err = (*c).CreateSecret(secret)
	default:
		return fmt.Errorf("unknown kind: %s", gvk.Kind)
	}

	return err
}
