package core

import (
	"interingo/pkg/service/common"
	"log/slog"
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
		slog.Error("API error, evalCore didn't init yet")
		return EvaluateResponse{Success: nil, Error: common.NewErrorResponse(500), Verbose: nil}
	}

	return c.evaluateHandler(*runtimeCore, req.Data)
}
