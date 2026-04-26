package v1

import (
	"fmt"
	"interingo/pkg/service/common"
	"interingo/pkg/service/core"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API v1 api/evalutate
func (s *interingoServiceV1Impl) EvaluateHandler(c *gin.Context) {
	req, err := evaluateHandlerRequest(c)
	if err != nil {
		slog.Error("API error, can't parse JSON value", "error", err.Error())
		errorResp := common.NewInvalidParamsErrorResponse(err.Error(), nil)
		evaluateHandlerErrorResponse(c, errorResp)
		return
	}

	if s.serviceCore == nil {
		slog.Error("API error, serviceCore didn't init yet")
		errorResp := common.NewErrorResponse(500)
		evaluateHandlerErrorResponse(c, errorResp)
		return
	}

	coreResp := s.serviceCore.EvaluateHandler(*req)
	if coreResp.Error != nil {
		evaluateHandlerErrorResponse(c, coreResp.Error)
		return
	}

	evaluateHandlerResponse(c, coreResp)
}

// Ensure craft user input to service core standard
func evaluateHandlerRequest(c *gin.Context) (*core.EvaluateRequest, error) {
	var userInput EvaluateRequest
	err := c.BindJSON(&userInput)

	if err != nil {
		slog.Error("API error, can't parse JSON value", "error", err.Error())
		return nil, fmt.Errorf("Can't parse JSON value")
	}
	req := &core.EvaluateRequest{
		Data: userInput.Data,
	}
	return req, nil
}

// Ensure output to api endpoit return to match API schema standard
func evaluateHandlerResponse(c *gin.Context, coreResp core.EvaluateResponse) {
	// Return
	resp := EvaluateResponseSuccess{
		Output: coreResp.Success.Output,
	}
	c.JSON(http.StatusOK, resp)
}

// Ensure output to api endpoit return to match API schema standard
func evaluateHandlerErrorResponse(c *gin.Context, error common.ErrorResponse) {
	switch error := error.(type) {
	case *core.EvalErrorResponse:
		resp := EvaluateResponseEvalError{
			Type:    error.GetType(),
			Code:    error.GetCode(),
			Message: error.GetMessage(),
		}
		c.JSON(resp.Type, resp)
	case *core.ParserErrorResponse:
		var parserError []ParserError
		for _, e := range error.Errors {
			parserError = append(parserError, ParserError{
				Message: e.Message,
			})
		}

		resp := EvaluateResponseParserError{
			Type:    error.GetType(),
			Code:    error.GetCode(),
			Message: error.GetMessage(),
			Errors:  parserError,
		}
		c.JSON(resp.Type, resp)
	default:
		resp := EvaluateResponseError{
			Type:    error.GetType(),
			Code:    error.GetCode(),
			Message: error.GetMessage(),
		}
		c.JSON(resp.Type, resp)
	}
}
