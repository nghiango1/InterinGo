package core

import "github.com/gorilla/websocket"

func (core *ServiceCore) WebsocketHandler(conn *websocket.Conn, messageType int, message []byte) {
	core.muConn.Lock()
	defer core.muConn.Unlock()

	if core.conn != nil {
		core.conn.WriteMessage(websocket.CloseMessage, []byte("Only support one ws client for print"))
		core.conn.Close()
	}

	core.conn = conn

	// if the client do want to try and close connection, remove core.conn
	// else we will clean up conn after print return error
	if messageType == websocket.CloseMessage {
		core.conn = nil
	}
}
