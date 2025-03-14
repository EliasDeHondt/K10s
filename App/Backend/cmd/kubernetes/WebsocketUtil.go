package kubernetes

import (
	"github.com/gorilla/websocket"
	"log"
)

func CloseConn(conn *websocket.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println("Error closing visualization socket:", err)
	}
}
