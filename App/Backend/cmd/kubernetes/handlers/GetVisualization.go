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

var conns = make(map[*websocket.Conn]string)

func GetVisualizationHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading visualization socket:", err)
		return
	}

	conns[conn] = ""

	VisualizationReady.Wait()
	cluster := CachedVisualization

	err = conn.WriteJSON(cluster)
	if wsError(err) {
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if wsError(err) {
			return
		} else {
			namespace := string(message)
			if namespace == "" {
				err = conn.WriteJSON(CachedVisualization)
			} else {
				err = conn.WriteJSON(CachedVisualization.FilterByNamespace(namespace))
			}
			if wsError(err) {
				return
			}
		}
		conns[conn] = string(message)
	}

}

func wsError(err error) bool {
	if err != nil {
		log.Println("Error writing visualization stats:", err)
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("WebSocket connection closed by client.")
		}
		return true
	}
	return false
}

func GetVisualizationConns() map[*websocket.Conn]string {
	return conns
}

func RemoveVisualizationConn(conn *websocket.Conn) {
	delete(conns, conn)
}
