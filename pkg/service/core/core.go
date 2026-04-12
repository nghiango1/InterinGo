package core

import (
	"fmt"
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"

	"github.com/gorilla/websocket"
)

type ServiceCore struct {
	runtimeCore *runtime.Core
	conn        *websocket.Conn
}

func NewServiceCore(evalCore *runtime.Core) *ServiceCore {
	if evalCore == nil {
		evalCore = runtime.NewCore(runtime.EMBED, nil)
	}

	res := &ServiceCore{
		runtimeCore: evalCore,
	}

	evalCore.Env.Set(
		"print", &PrintBuiltin{
			env:  res.runtimeCore.Env,
			core: res,
		},
	)

	return res
}

// Return
// Success -> EvaluateResponseSuccess
// Error -> 400:Bad request
// Error -> 500:Internal server error
func (c *ServiceCore) EvaluateHandler(req EvaluateRequest) (*EvaluateResponseSuccess, common.ErrorResponseInterface) {
	if c.runtimeCore == nil {
		fmt.Println("[ERRPR] API error, evalCore didn't init yet")
		return nil, common.NewErrorResponse(500)
	}

	// Handling eval
	res, err, ver := c.runtimeCore.Eval(runtime.EvalRequest{Data: req.Data})

	if err != nil {
		error := NewParserErrorResponse("", err.ParserErrors, ver)
		return nil, error
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
