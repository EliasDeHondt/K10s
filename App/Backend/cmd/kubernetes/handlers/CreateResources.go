/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"bytes"
	"fmt"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"strings"
)

func CreateResourcesHandler(ctx *gin.Context) {
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

	obj, err := CreateResources(C, data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, obj)
}

func CreateResources(c kubernetes.IClient, data []byte) ([]interface{}, error) {
	decoder := scheme.Codecs.UniversalDeserializer()

	yamlResources := bytes.Split(data, []byte("---"))

	var madeResources []interface{}

	for _, y := range yamlResources {
		if isCommentsOrWhitespaceOnly(y) {
			continue
		}

		obj, gvk, err := decoder.Decode(y, nil, nil)
		if err != nil {
			return madeResources, err
		}

		var newResource interface{}

		switch gvk.Kind {
		case "Namespace":
			namespace := obj.(*corev1.Namespace)
			newResource, err = c.CreateNamespace(namespace)
			madeResources = append(madeResources, newResource)
		case "Node":
			node := obj.(*corev1.Node)
			newResource, err = c.CreateNode(node)
			madeResources = append(madeResources, newResource)
		case "Pod":
			pod := obj.(*corev1.Pod)
			newResource, err = c.CreatePod(pod)
			madeResources = append(madeResources, newResource)
		case "Deployment":
			deployment := obj.(*appsv1.Deployment)
			newResource, err = c.CreateDeployment(deployment)
			madeResources = append(madeResources, newResource)
		case "Service":
			service := obj.(*corev1.Service)
			var newService kubernetes.Service
			newService, err = c.CreateService(service)
			madeResources = append(madeResources, newService)
		case "ConfigMap":
			configMap := obj.(*corev1.ConfigMap)
			newResource, err = c.CreateConfigMap(configMap)
			madeResources = append(madeResources, newResource)
		case "Secret":
			secret := obj.(*corev1.Secret)
			newResource, err = c.CreateSecret(secret)
			madeResources = append(madeResources, newResource)
		}

		if err != nil {
			return madeResources, fmt.Errorf("failed to apply %s: %w", gvk.Kind, err)
		}

	}

	return madeResources, nil
}

func isCommentsOrWhitespaceOnly(data []byte) bool {
	lines := bytes.Split(data, []byte("\n"))
	whitespaceCount := 0
	hasCommentedLinesOnly := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(string(line))
		if trimmed == "" {
			whitespaceCount++
		} else {
			if !strings.HasPrefix(trimmed, "#") {
				hasCommentedLinesOnly = false
			}
		}
	}

	return (whitespaceCount == len(lines)) || hasCommentedLinesOnly

}
