package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == GetFrontendIP()
	},
}

func HandleMetricsSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading metrics socket:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing metrics socket:", err)
		}
	}(conn)

	for {
		stats, err := c.GetTotalUsage()
		if err != nil {
			log.Println("Error getting metrics stats:", err)
			return
		}

		err = conn.WriteJSON(stats)
		if err != nil {
			log.Println("Error writing metrics stats:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket connection closed by client.")
			}
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}
