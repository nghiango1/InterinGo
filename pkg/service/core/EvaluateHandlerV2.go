package core

import (
	"fmt"
	"interingo/pkg/service/common"
)

func (c *ServiceCore) EvaluateHandlerV2(req EvaluateRequest) EvaluateResponse {
	if req.RuntimeId == "" {
		return c.EvaluateHandler(req)
	}
	runtimeCore, ok := c.runtimeCores[req.RuntimeId]
	if !ok {
		return EvaluateResponse{Success: nil, Error: common.NewErrorResponse(404), Verbose: nil}
	}

	if runtimeCore == nil {
		fmt.Println("[ERROR] API error, evalCore didn't init yet")
		return EvaluateResponse{Success: nil, Error: common.NewErrorResponse(500), Verbose: nil}
	}

	return c.evaluateHandler(*runtimeCore, req.Data)
}
