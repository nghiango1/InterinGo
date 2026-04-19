package core

import (
	"fmt"
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"
)

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
