package core

import (
	"fmt"
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectedClient struct {
	muConn sync.Mutex
	conn   *websocket.Conn
}

type ReplRuntime struct {
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

// Return
// Success -> EvaluateResponseSuccess
// Error -> 400:Bad request
// Error -> 500:Internal server error
func (c *ServiceCore) EvaluateHandler(req EvaluateRequest) (*EvaluateResponseSuccess, common.ErrorResponseInterface) {
	if c.runtimeCoreDefault == nil {
		fmt.Println("[ERRPR] API error, evalCore didn't init yet")
		return nil, common.NewErrorResponse(500)
	}

	res, err, ver := c.runtimeCoreDefault.Eval(runtime.EvalRequest{Data: req.Data})

	if err != nil {
		message := ""
		if err.Error != nil {
			message = err.Error.Error()
		}

		if err.SystemExit != nil {
			return &EvaluateResponseSuccess{
				Output:  nil,
				Verbose: ver,
			}, nil
		}

		if len(err.ParserErrors) != 0 {
			error := NewParserErrorResponse(message, err.ParserErrors, ver)
			return nil, error
		}

		if err.Error != nil {
			return nil, NewEvalErrorResponse(message, ver)
		}

		return nil, common.NewErrorResponse(500) // Unknown error
	} else if res != nil {
		success := EvaluateResponseSuccess{
			Output:  res.Output,
			Verbose: ver,
		}
		return &success, nil
	}

	// If both err and res from Eval is nil, there some thing wrong
	return nil, common.NewErrorResponse(500)
}

// Return
// Success -> EvaluateResponseSuccess
// Error -> 400:Bad request
// Error -> 500:Internal server error
func (c *ServiceCore) EvaluateHandlerV2(req EvaluateRequest) (*EvaluateResponseSuccess, common.ErrorResponseInterface) {
	if req.RuntimeId == "" {
		return c.EvaluateHandler(req)
	}
	runtimeCore, ok := c.runtimeCores[req.RuntimeId]
	if !ok {
		return nil, common.NewErrorResponse(404)
	}

	if runtimeCore != nil {
		fmt.Println("[ERROR] API error, evalCore didn't init yet")
		return nil, common.NewErrorResponse(500)
	}

	// Handling eval can cause print, which will use connClients
	connClient, ok := c.connClients[runtimeCore.connId]
	if ok {
		connClient.muConn.Lock()
		defer connClient.muConn.Unlock()
	}
	res, err, ver := runtimeCore.core.Eval(runtime.EvalRequest{Data: req.Data})

	if err != nil {
		message := ""
		if err.Error != nil {
			message = err.Error.Error()
		}
		if err.SystemExit != nil {
			return &EvaluateResponseSuccess{
				Output:  nil,
				Verbose: ver,
			}, nil
		}

		if len(err.ParserErrors) != 0 {
			error := NewParserErrorResponse(message, err.ParserErrors, ver)
			return nil, error
		}

		if err.Error != nil {
			return nil, NewEvalErrorResponse(message, ver)
		}

		{
			return nil, common.NewErrorResponse(500) // Unknown error
		}
	} else if res != nil {
		success := EvaluateResponseSuccess{
			Output:  res.Output,
			Verbose: ver,
		}
		return &success, nil
	}

	// If both err and res from Eval is nil, there some thing wrong
	return nil, common.NewErrorResponse(500)
}
