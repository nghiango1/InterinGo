package core

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (core *ServiceCore) WebsocketConnectionCreate(conn *websocket.Conn) (string, error) {
	core.muConnClients.Lock()
	defer core.muConnClients.Unlock()

	connId := uuid.New().String()
	_, ok := core.connClients[connId]
	if ok {
		log.Printf("[ERROR] ConnId collision, should not be possible")
		return "", fmt.Errorf("[ERROR] ConnId collision, should not be possible")
	}

	client := &ConnectedClient{}
	client.muConn.Lock()
	defer client.muConn.Unlock()

	client.conn = conn
	core.connClients[connId] = client

	log.Printf("[INFO] New connection: %v", NewWebsocketConnectSuccess(connId))
	conn.WriteJSON(NewWebsocketConnectSuccess(connId))

	return connId, nil
}

func (core *ServiceCore) WebsocketConnectionCleanup(connId string) {
	core.muConnClients.Lock()
	defer core.muConnClients.Unlock()

	delete(core.connClients, connId)
}
