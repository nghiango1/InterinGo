package core

import "github.com/gorilla/websocket"

func (core *ServiceCore) WebsocketHandler(conn *websocket.Conn, messageType int, message []byte) {
	// Ignore every thing client have to said
}
