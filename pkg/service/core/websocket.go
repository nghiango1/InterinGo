package core

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (core *ServiceCore) WebsocketConnectionCreate(conn *websocket.Conn) string {
	core.muConn.Lock()
	defer core.muConn.Unlock()

	if core.conn != nil {
		core.conn.WriteMessage(websocket.CloseMessage, []byte("Cuurrently only support one ws client for print"))
		core.conn.Close()
	}

	core.conn = conn
	core.connId = uuid.New().String()

	return core.connId
}

func (core *ServiceCore) WebsocketHandler(connId string, messageType int, message []byte) error {
	core.muConn.Lock()
	defer core.muConn.Unlock()

	if core.connId != connId {
		return fmt.Errorf("ConnId doesn't match, server should drop the connection")
	}

	// if the client do want to try and close connection, remove core.conn
	// else we will clean up conn after print return error
	if messageType == websocket.CloseMessage {
		core.conn = nil
	}
	return nil
}
