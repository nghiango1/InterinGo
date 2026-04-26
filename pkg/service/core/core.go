package core

import (
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	alivePeriod = 60 * time.Second
)

type ConnectedClient struct {
	id      string
	conn    *websocket.Conn
	runtime *ReplRuntime
}

func NewConnectedClient(conn *websocket.Conn) *ConnectedClient {
	connId := uuid.New().String()

	client := &ConnectedClient{
		id:      connId,
		conn:    conn,
		runtime: nil,
	}

	return client
}

type ReplRuntime struct {
	id         string
	core       *runtime.Core
	connId     string
	lastUpdate time.Time
}

type ServiceCore struct {
	muConnClients      sync.Mutex
	connClients        map[string]*ConnectedClient
	runtimeCores       map[string]*ReplRuntime
	runtimeCoreDefault *runtime.Core
}

func (core *ServiceCore) EnsureDefaultEvalCode() {
	if core.runtimeCoreDefault != nil {
		if core.runtimeCoreDefault.State() != runtime.RUNTIME_END {
			return
		}
	}

	evalCore := runtime.NewCore(runtime.EMBED, nil)

	evalCore.Env.Set(
		"print", &PrintBuiltin{
			env: evalCore.Env,
			externalPrint: func(message string) {
				slog.Debug("Global PrintBuiltin request lock muConnClients")
				core.muConnClients.Lock()
				slog.Debug("Global PrintBuiltin take lock muConnClients")
				defer core.muConnClients.Unlock()
				defer slog.Debug("Global PrintBuiltin release lock muConnClients")

				for _, client := range core.connClients {
					client.conn.WriteJSON(NewPrintMessageEventData(message))
				}
			},
		},
	)

	core.runtimeCoreDefault = evalCore
}

func NewServiceCore(evalCore *runtime.Core) *ServiceCore {
	serviceCore := &ServiceCore{
		connClients:        make(map[string]*ConnectedClient),
		runtimeCores:       make(map[string]*ReplRuntime),
		runtimeCoreDefault: evalCore,
	}

	if evalCore == nil {
		serviceCore.EnsureDefaultEvalCode()
	}

	// With out Client connection, we need to clean up base from the request
	// for evaluate of this REPL session
	go func() {
		// good enough, timer check every 60 second
		ticker := time.NewTicker(alivePeriod)
		defer ticker.Stop()
		for range ticker.C {
			slog.Debug("Service Core Clean up request lock muConnClients")
			serviceCore.muConnClients.Lock()
			slog.Debug("Service Core Clean up take lock muConnClients")

			serviceCore.EnsureDefaultEvalCode()
			// else check if lastUpdate + alivePeriod < current time.
			// or there already have a client bind to it
			// If so clean up is perform
			for _, runtime := range serviceCore.runtimeCores {
				if time.Now().After(runtime.lastUpdate.Add(alivePeriod)) && runtime.connId == "" {
					slog.Debug("Service core release runtime", "runtime_id", runtime.id)
					delete(serviceCore.runtimeCores, runtime.id)
				}
			}

			serviceCore.muConnClients.Unlock()
			slog.Debug("Service Core Clean up release lock muConnClients")
		}
	}()

	return serviceCore
}

// common handler
func (c *ServiceCore) evaluateHandler(runtimeCore ReplRuntime, data string) EvaluateResponse {
	res, err, ver := runtimeCore.core.Eval(runtime.EvalRequest{Data: data})
	runtimeCore.lastUpdate = time.Now()

	if err != nil {
		message := ""
		if err.Error != nil {
			message = err.Error.Error()
		}

		if err.SystemExit != nil {
			return EvaluateResponse{
				Success: &EvaluateResponseSuccess{
					Output: nil,
				},
				Error:   nil,
				Verbose: ver,
			}
		}

		if len(err.ParserErrors) != 0 {
			error := NewParserErrorResponse(message, err.ParserErrors)
			return EvaluateResponse{
				Success: nil,
				Error:   error,
				Verbose: ver,
			}
		}

		if err.Error != nil {
			return EvaluateResponse{
				Success: nil,
				Error:   NewEvalErrorResponse(message),
				Verbose: ver,
			}
		}

		return EvaluateResponse{
			Success: nil,
			Error:   common.NewErrorResponse(500),
			Verbose: ver,
		} // Unknown error
	} else if res != nil {
		success := EvaluateResponseSuccess{
			Output: res.Output,
		}
		return EvaluateResponse{
			Success: &success,
			Error:   nil,
			Verbose: ver,
		}
	}

	// If both err and res from Eval is nil, there some thing wrong
	return EvaluateResponse{
		Success: nil,
		Error:   common.NewErrorResponse(500),
		Verbose: ver,
	}
}
