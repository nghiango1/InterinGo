package core

import (
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"
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

// common handler
func (c *ServiceCore) evaluateHandler(runtimeCore ReplRuntime, data string) EvaluateResponse {
	res, err, ver := runtimeCore.core.Eval(runtime.EvalRequest{Data: data})

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
