/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package handlers

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func GetVisualizationHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading visualization socket:", err)
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing visualization socket:", err)
		}
	}(conn)

	for {
		kubernetes.VisualizationReady.Wait()
		cluster := kubernetes.CachedVisualization

		err = conn.WriteJSON(cluster)
		if err != nil {
			log.Println("Error writing visualization stats:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket connection closed by client.")
			}
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}
