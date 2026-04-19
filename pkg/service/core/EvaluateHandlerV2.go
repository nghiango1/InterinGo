package core

import (
	"fmt"
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"
)

func (c *ServiceCore) EvaluateHandlerV2(req EvaluateRequest) (*EvaluateResponseSuccess, common.ErrorResponseInterface) {
	if req.RuntimeId == "" {
		return c.EvaluateHandler(req)
	}
	runtimeCore, ok := c.runtimeCores[req.RuntimeId]
	if !ok {
		return nil, common.NewErrorResponse(404)
	}

	if runtimeCore == nil {
		fmt.Println("[ERROR] API error, evalCore didn't init yet")
		return nil, common.NewErrorResponse(500)
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
