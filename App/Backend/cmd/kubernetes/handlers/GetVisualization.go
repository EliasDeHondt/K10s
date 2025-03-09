package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
)

func GetVisualizationHandler(ctx *gin.Context) {
	cluster := kubernetes.VisualizeCluster(*c)
	ctx.JSON(200, cluster)
}
