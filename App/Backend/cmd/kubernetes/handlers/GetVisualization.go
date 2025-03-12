package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
)

func GetVisualizationHandler(ctx *gin.Context) {
	namespace, _ := ctx.GetQuery("namespace")
	cluster := kubernetes.VisualizeCluster(c, namespace)
	ctx.JSON(200, cluster)
}
