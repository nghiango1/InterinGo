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

type Message interface {
	Type() MessageType
}

type ReplBindMessage struct {
	MessageType MessageType `json:"type"`
	RuntimeId   string      `json:"runtimeId"`
}

func (m *ReplBindMessage) Type() MessageType {
	return m.MessageType
}

type MessageType string

const (
	REPL_BIND = MessageType("repl_bind")
)

func (core *ServiceCore) WebsocketMessageHandler(connectedClient *ConnectedClient, data []byte) error {
	log.Printf("[INFO] ServiceCore Websocket handler got request %v", string(data))
	var mes ReplBindMessage
	err := json.Unmarshal(data, &mes)

	if err != nil {
		log.Printf("[ERROR] Failed to read message data")
		return err
	}

	return core.WebsocketConnectionReplBind(connectedClient, mes.RuntimeId)
}

func (core *ServiceCore) WebsocketConnectionReplBind(connectedClient *ConnectedClient, runtimeId string) error {
	core.muConnClients.Lock()
	defer core.muConnClients.Unlock()

	runtime, ok := core.runtimeCores[runtimeId]
	if !ok {
		return fmt.Errorf("[ERROR] Runtime not found")
	}

	// You can try to connect to others people created runtime, which overide
	// the connection and then it can't received printed message anymore

	// Clean up runtimeId which this connectedClient connected to
	// UI only allow one REPL anyway, so this is safe
	if connectedClient.runtimeId != "" {
		prevRuntime, ok := core.runtimeCores[connectedClient.runtimeId]
		if ok && prevRuntime != nil {
			delete(core.runtimeCores, prevRuntime.id)
		}
	}

	runtime.core.Env.Set(
		"print", &PrintBuiltin{
			env: runtime.core.Env,
			externalPrint: func(message string) {
				connectedClient.muConn.Lock()
				defer connectedClient.muConn.Unlock()

				connectedClient.conn.WriteJSON(NewPrintMessageEventData(message))
			},
		},
	)

	// Overide to the new runtime
	connectedClient.runtimeId = runtimeId
	runtime.connId = connectedClient.id
	return nil
}

func (core *ServiceCore) WebsocketConnectionCleanup(connectedClient *ConnectedClient) {
	core.muConnClients.Lock()
	defer core.muConnClients.Unlock()

	// Clean up runtimeId which this connectedClient connected to
	// UI only allow one REPL anyway, so this is safe
	if connectedClient.runtimeId != "" {
		prevRuntime, ok := core.runtimeCores[connectedClient.runtimeId]
		if ok && prevRuntime != nil {
			delete(core.runtimeCores, prevRuntime.id)
		}
	}

	delete(core.connClients, connectedClient.id)
}
