/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package main

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/auth"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

var frontendUrl = handlers.GetFrontendIP()

func main() {
	frontendUrl = handlers.GetFrontendIP()
	//trustedProxies := []string{"10.0.0.0/8"}

	if frontendUrl == "http://localhost:4200" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth.Init()
	r.POST("/api/login", auth.HandleLogin)
	r.GET("/api/logout", auth.HandleLogout)
	r.GET("/api/isloggedin", auth.IsLoggedIn)

	secured := r.Group("/api/secured")
	secured.Use(auth.AuthMiddleware())
	secured.GET("/table", handlers.GetTableHandler)
	secured.GET("/nodes", handlers.GetNodesHandler)
	secured.GET("/pods", handlers.GetPodsHandler)
	secured.GET("/services", handlers.GetServicesHandler)
	secured.GET("/configMaps", handlers.GetConfigMapsHandler)
	secured.GET("/secrets", handlers.GetSecretsHandler)
	secured.GET("/deployments", handlers.GetDeploymentsHandler)
	secured.POST("/createresources", handlers.CreateResourcesHandler)
	secured.GET("/namespaces", handlers.GetNamespacesHandler)
	secured.GET("/nodenames", handlers.GetNodeNamesHandler)
	secured.GET("/statsocket", handlers.HandleMetricsSocket)
	secured.GET("/visualization", handlers.GetVisualizationHandler)

	handlers.VisualizationReady.Add(1)
	go handlers.CreateVisualization(handlers.C)
	go handlers.C.WatchUsage()

	err = r.Run(":8082")
	if err != nil {
		panic(err)
	}
}
