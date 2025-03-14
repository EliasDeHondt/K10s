/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

var conns = make([]*websocket.Conn, 0)

func GetVisualizationHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading visualization socket:", err)
		return
	}

	conns = append(conns, conn)

	VisualizationReady.Wait()
	cluster := CachedVisualization

	err = conn.WriteJSON(cluster)

	if err != nil {
		log.Println("Error writing visualization stats:", err)
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("WebSocket connection closed by client.")
			return
		}
	}
}

func GetVisualizationConns() []*websocket.Conn {
	return conns
}
