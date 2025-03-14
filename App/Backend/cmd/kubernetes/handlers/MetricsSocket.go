/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func HandleMetricsSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading metrics socket:", err)
		return
	}

	C.AddMetricsConnection(conn)

}
