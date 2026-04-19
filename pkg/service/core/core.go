package core

import (
	"interingo/pkg/runtime"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectedClient struct {
	muConn sync.Mutex
	conn   *websocket.Conn
	id     string
}

type ReplRuntime struct {
	id     string
	core   *runtime.Core
	connId string
}

type ServiceCore struct {
	muConnClients      sync.Mutex
	connClients        map[string]*ConnectedClient
	runtimeCores       map[string]*ReplRuntime
	runtimeCoreDefault *runtime.Core
}

func NewServiceCore(evalCore *runtime.Core) *ServiceCore {
	evalCore = runtime.NewCore(runtime.EMBED, nil)

	serviceCore := &ServiceCore{
		connClients:        make(map[string]*ConnectedClient),
		runtimeCores:       make(map[string]*ReplRuntime),
		runtimeCoreDefault: evalCore,
	}

	evalCore.Env.Set(
		"print", &PrintBuiltin{
			env: evalCore.Env,
			externalPrint: func(message string) {
				serviceCore.muConnClients.Lock()
				defer serviceCore.muConnClients.Unlock()

				for _, client := range serviceCore.connClients {
					client.conn.WriteJSON(NewPrintMessageEventData(message))
				}
			},
		},
	)

	return serviceCore
}
