/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Element string

const (
	Nodes       Element = "nodes"
	Pods        Element = "pods"
	Services    Element = "services"
	Deployments Element = "deployments"
	ConfigMaps  Element = "configmaps"
	Secrets     Element = "secrets"
)

func isValidElement(e Element) bool {
	switch e {
	case Nodes, Pods, Services, Deployments, ConfigMaps, Secrets:
		return true
	default:
		return false
	}
}

func GetTableHandler(ctx *gin.Context) {

	elementType, elementOk := ctx.GetQuery("element")

	if !elementOk {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Missing query parameter 'element'"})
		return
	}

	element := Element(strings.ToLower(elementType))
	if !isValidElement(element) {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid query parameter 'element'"})
		return
	}

	switch element {
	case Nodes:
		GetNodesHandler(ctx)
	case Pods:
		GetPodsHandler(ctx)
	case Services:
		GetServicesHandler(ctx)
	case Deployments:
		GetDeploymentsHandler(ctx)
	case ConfigMaps:
		GetConfigMapsHandler(ctx)
	case Secrets:
		GetSecretsHandler(ctx)
	}
}