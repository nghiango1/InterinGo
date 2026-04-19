package v2

import (
	"fmt"
	"interingo/pkg/service/common"
	"interingo/pkg/service/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API v2 api/repl/:id/evalutate
func (s *interingoServiceV2Impl) EvaluateHandler(c *gin.Context) {
	req, err := evaluateHandlerRequest(c)
	if err != nil {
		fmt.Println("[ERRPR] API error, can't parse JSON value, got: ", err.Error())
		errorResp := common.NewBadRequestErrorResponse(err.Error(), nil)
		evaluateHandlerErrorResponse(c, errorResp)
		return
	}

	if s.serviceCore == nil {
		fmt.Println("[ERROR] API error, serviceCore didn't init yet")
		errorResp := common.NewErrorResponse(500)
		evaluateHandlerErrorResponse(c, errorResp)
		return
	}

	resp := s.serviceCore.EvaluateHandlerV2(*req)

	// Return
	if resp.Error != nil {
		c.JSON(resp.Error.GetType(), resp.Error)
	} else if resp.Success != nil {
		c.JSON(http.StatusOK, resp.Success)
	}
}

// Ensure craft user input to service core standard
func evaluateHandlerRequest(c *gin.Context) (*core.EvaluateRequest, error) {
	runtimeId := c.Param("id")
	// Input validate
	var userInput EvaluateRequest
	err := c.BindJSON(&userInput)
	if err != nil {
		fmt.Println("[ERROR] API error, can't parse JSON value, got: ", err.Error())
		return nil, fmt.Errorf("Can't parse JSON value")
	}

	req := &core.EvaluateRequest{
		RuntimeId: runtimeId,
		Data:      userInput.Data,
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
func evaluateHandlerErrorResponse(c *gin.Context, error common.ErrorResponseInterface) {
	switch error := error.(type) {
	case *core.EvalErrorResponse:
		resp := EvaluateResponseEvalError{
			Type:    error.GetType(),
			Code:    error.GetCode(),
			Message: error.GetMessage(),
		}
		c.JSON(resp.Type, resp)
	case *core.ParserErrorResponse:
		resp := EvaluateResponseParserError{
			Type:    error.GetType(),
			Code:    error.GetCode(),
			Message: error.GetMessage(),
			Errors:  error.Errors,
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
