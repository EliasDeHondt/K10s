/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package main

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/auth"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	auth.Init()
	r.POST("/login", auth.HandleLogin)
	r.GET("/logout", auth.HandleLogout)

	secured := r.Group("/secured")
	secured.Use(auth.AuthMiddleware())
	secured.GET("/table", handlers.GetTableHandler)
	secured.GET("/nodes", handlers.GetNodesHandler)
	secured.GET("/pods", handlers.GetPodsHandler)
	secured.GET("/services", handlers.GetServicesHandler)
	secured.GET("/configMaps", handlers.GetConfigMapsHandler)
	secured.GET("/secrets", handlers.GetSecretsHandler)
	secured.GET("/deployments", handlers.GetDeploymentsHandler)
	secured.GET("/stats", handlers.GetStatsHandler)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
