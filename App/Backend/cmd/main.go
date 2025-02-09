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

	r.POST("/login", auth.HandleLogin)
	r.GET("/logout", auth.HandleLogout)

	secured := r.Group("/secured")
	secured.Use(auth.AuthMiddleware())
	//TODO: Move to secured after testing
	r.GET("/nodes", handlers.GetNodesHandler)
	r.GET("/pods", handlers.GetPodsHandler)
	r.GET("/services", handlers.GetServicesHandler)
	r.GET("/configMaps", handlers.GetConfigMapsHandler)
	r.GET("/secrets", handlers.GetSecretsHandler)
	r.GET("/deployments", handlers.GetDeploymentsHandler)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
