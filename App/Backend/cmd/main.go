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
	r.GET("/nodes", handlers.GetNodesHandler)
	r.GET("/pods", handlers.GetPodsHandler)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
