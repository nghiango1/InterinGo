package core

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func (core *ServiceCore) WebsocketConnectionCreate(conn *websocket.Conn) (*ConnectedClient, error) {
	core.muConnClients.Lock()
	defer core.muConnClients.Unlock()

	connectedClient := NewConnectedClient(conn)
	_, ok := core.connClients[connectedClient.id]
	if ok {
		log.Printf("[ERROR] ConnId collision, should not be possible")
		return nil, fmt.Errorf("[ERROR] ConnId collision, should not be possible")
	}

	core.connClients[connectedClient.id] = connectedClient

	log.Printf("[INFO] New connection: %v", NewWebsocketConnectSuccess(connectedClient.id))
	conn.WriteJSON(NewWebsocketConnectSuccess(connectedClient.id))

	return connectedClient, nil
}

func (core *ServiceCore) WebsocketReceivedTextMessageHandler(connectedClient *ConnectedClient, data []byte) error {
	log.Printf("[INFO] ServiceCore Websocket handler got request %v", string(data))
	var mes ReplBindRequest
	err := json.Unmarshal(data, &mes)

	if err != nil {
		log.Printf("[ERROR] Failed to read message data")
		return err
	}

	return core.WebsocketReplBindHandler(connectedClient, mes.RuntimeId)
}

func (core *ServiceCore) WebsocketReplBindHandler(connectedClient *ConnectedClient, runtimeId string) error {
	core.muConnClients.Lock()
	defer core.muConnClients.Unlock()

	runtime, ok := core.runtimeCores[runtimeId]
	if !ok {
		return fmt.Errorf("[ERROR] Runtime not found")
	}

	// You can try to connect to others people created runtime, which overide
	// the connection and then it can't received printed message anymore

	runtime.core.Env.Set(
		"print", &PrintBuiltin{
			env: runtime.core.Env,
			externalPrint: func(message string) {
				connectedClient.conn.WriteJSON(NewPrintMessageEventData(message))
			},
		},
	)

	if connectedClient.runtime != nil {
		log.Printf("[INFO] Connection %s release runtime %s", connectedClient.id, runtime.id)
	}
	// Overide to the new runtime
	connectedClient.runtime = runtime
	// Remove binding to this runtime from core (connectedClient own it now)
	delete(core.runtimeCores, runtime.id)
	log.Printf("[INFO] Connection %s take over runtime %s", connectedClient.id, runtime.id)

	runtime.connId = connectedClient.id
	return nil
}

func (core *ServiceCore) WebsocketConnectionCleanup(connectedClient *ConnectedClient) {
	core.muConnClients.Lock()
	defer core.muConnClients.Unlock()

	delete(core.connClients, connectedClient.id)
}
