package core

import (
	"fmt"
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"
)

func (c *ServiceCore) EvaluateHandlerV2(req EvaluateRequest) EvaluateResponse {
	if req.RuntimeId == "" {
		return c.EvaluateHandler(req)
	}
	runtimeCore, ok := c.runtimeCores[req.RuntimeId]
	if !ok {
		return EvaluateResponse{Success: nil, Error: common.NewErrorResponse(404)}
	}

	if runtimeCore == nil {
		fmt.Println("[ERROR] API error, evalCore didn't init yet")
		return EvaluateResponse{Success: nil, Error: common.NewErrorResponse(500)}
	}

	res, err, ver := runtimeCore.core.Eval(runtime.EvalRequest{Data: req.Data})

	if err != nil {
		message := ""
		if err.Error != nil {
			message = err.Error.Error()
		}
		if err.SystemExit != nil {
			return EvaluateResponse{Success: &EvaluateResponseSuccess{
				Output:  nil,
				Verbose: ver,
			}, Error: nil}
		}

		if len(err.ParserErrors) != 0 {
			error := NewParserErrorResponse(message, err.ParserErrors, ver)
			return EvaluateResponse{Success: nil, Error: error}
		}

		if err.Error != nil {
			return EvaluateResponse{Success: nil, Error: NewEvalErrorResponse(message, ver)}
		}

		return EvaluateResponse{Success: nil, Error: common.NewErrorResponse(500)} // Unknown error
	} else if res != nil {
		success := EvaluateResponseSuccess{
			Output:  res.Output,
			Verbose: ver,
		}
		return EvaluateResponse{Success: &success, Error: nil}
	}

	// If both err and res from Eval is nil, there some thing wrong
	return EvaluateResponse{Success: nil, Error: common.NewErrorResponse(500)}
}
