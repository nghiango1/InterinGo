package core

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

func (core *ServiceCore) WebsocketConnectionCreate(conn *websocket.Conn) (*ConnectedClient, error) {
	slog.Debug("WebsocketConnectionCreate request lock muConnClients")
	core.muConnClients.Lock()
	slog.Debug("WebsocketConnectionCreate take lock muConnClients")
	defer core.muConnClients.Unlock()
	defer slog.Debug("WebsocketConnectionCreate release lock muConnClients")

	connectedClient := NewConnectedClient(conn)
	_, ok := core.connClients[connectedClient.id]
	if ok {
		slog.Error("ConnId collision, should not be possible")
		return nil, fmt.Errorf("[ERROR] ConnId collision, should not be possible")
	}

	core.connClients[connectedClient.id] = connectedClient

	slog.Info("New websocket connection", "connection_id", connectedClient.id)
	conn.WriteJSON(NewWebsocketConnectSuccess(connectedClient.id))

	return connectedClient, nil
}

func (core *ServiceCore) WebsocketReceivedTextMessageHandler(connectedClient *ConnectedClient, data []byte) error {
	slog.Debug("ServiceCore Websocket handler got request", "data", string(data))
	var mes ReplBindRequest
	err := json.Unmarshal(data, &mes)

	if err != nil {
		slog.Error("Failed to read message data", "error", err)
		return err
	}

	return core.WebsocketReplBindHandler(connectedClient, mes.RuntimeId)
}

func (core *ServiceCore) WebsocketReplBindHandler(connectedClient *ConnectedClient, runtimeId string) error {
	slog.Debug("WebsocketReplBindHandler request lock muConnClients")
	core.muConnClients.Lock()
	slog.Debug("WebsocketReplBindHandler take lock muConnClients")
	defer core.muConnClients.Unlock()
	defer slog.Debug("WebsocketReplBindHandler release lock muConnClients")

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
		slog.Info("Connection release runtime", "connection_id", connectedClient.id, "runtime_id", runtime.id)
	}
	// Overide to the new runtime
	connectedClient.runtime = runtime

	runtime.connId = connectedClient.id
	return nil
}

func (core *ServiceCore) WebsocketConnectionCleanup(connectedClient *ConnectedClient) {
	slog.Debug("WebsocketConnectionCleanup request lock muConnClients")
	core.muConnClients.Lock()
	slog.Debug("WebsocketConnectionCleanup take lock muConnClients")
	defer core.muConnClients.Unlock()
	defer slog.Debug("WebsocketConnectionCleanup release lock muConnClients")

	// Still have some room for others to take over the Repl Runtime
	// till the next Core Cleanup cycle
	delete(core.connClients, connectedClient.id)
	connectedClient.runtime = nil
}
