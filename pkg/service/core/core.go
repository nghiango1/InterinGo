package core

import (
	"fmt"
	"interingo/pkg/evaluator"
	"interingo/pkg/service/common"
)

type ServiceCore struct {
	evalCore *evaluator.Core
}

func NewServiceCore(evalCore *evaluator.Core) *ServiceCore {
	if evalCore == nil {
		evalCore = evaluator.NewCore()
	}

	return &ServiceCore{
		evalCore: evalCore,
	}
}

// Return
// Success -> EvaluateResponseSuccess
// Error -> 400:Bad request 
// Error -> 500:Internal server error
func (c *ServiceCore) EvaluateHandler(req EvaluateRequest) (*EvaluateResponseSuccess, common.ErrorResponseInterface) {
	if c.evalCore == nil{
		fmt.Println("[ERRPR] API error, evalCore didn't init yet")
		return nil, common.NewErrorResponse(500)
	}

	// Handling eval
	res, err, ver := c.evalCore.Eval(evaluator.EvalRequest{Data: req.Data})

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
