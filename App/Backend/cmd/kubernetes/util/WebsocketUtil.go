package util

import (
	"github.com/gorilla/websocket"
	"log"
)

func CloseConn(conn *websocket.Conn, source string) {
	err := conn.Close()
	if err != nil {
		log.Printf("Error closing %s socket: %s \n", source, err)
	}
}
