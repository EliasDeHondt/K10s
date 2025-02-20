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
	"os"
	"time"
)

var envFrontendUrl, frontendUrlOk = os.LookupEnv("FRONTEND_URL")

func main() {

	if !frontendUrlOk {
		envFrontendUrl = "http://localhost:8080"
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{envFrontendUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth.Init()
	r.POST("/login", auth.HandleLogin)
	r.GET("/logout", auth.HandleLogout)
	r.GET("/isloggedin", auth.IsLoggedIn)

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
	secured.POST("/createresources", handlers.CreateResourcesHandler)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
