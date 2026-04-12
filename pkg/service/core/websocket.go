package core

import "github.com/gorilla/websocket"

func (core *ServiceCore) WebsocketHandler(conn *websocket.Conn, messageType int, message []byte) {
	if core.conn != nil {
		core.conn.WriteMessage(websocket.CloseMessage, []byte("Only support one ws client for print"))
		core.conn.Close()
	}
	core.conn = conn
	// Ignore every thing client have to said
}
