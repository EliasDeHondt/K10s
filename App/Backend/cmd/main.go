/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package main

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/auth"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.POST("/login", auth.HandleLogin)
	r.GET("/logout", auth.HandleLogout)

	secured := r.Group("/secured")
	secured.Use(auth.AuthMiddleware())
	secured.GET("/", func(c *gin.Context) {
		println("Test print")
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
