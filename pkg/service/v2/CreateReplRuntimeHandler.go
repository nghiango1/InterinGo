package v2

import (
	"fmt"
	"interingo/pkg/service/common"
	"interingo/pkg/service/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API v2 POST api/repl
func (s *interingoServiceV2Impl) CreateReplRuntimeHandler(c *gin.Context) {
	// Input validate
	req, err := createReplRuntimeRequest(c)
	if err != nil {
		fmt.Println("[ERRPR] API error, can't parse JSON value, got: ", err.Error())
		errorResp := common.NewInvalidParamsErrorResponse(err.Error(), nil)
		createReplRuntimeErrorResponse(c, errorResp)
		return
	}

	if s.serviceCore == nil {
		fmt.Println("[ERRPR] API error, serviceCore didn't init yet")
		errorResp := common.NewErrorResponse(500)
		createReplRuntimeErrorResponse(c, errorResp)
		return
	}

	coreResp, coreErr := s.serviceCore.CreateReplRuntime(*req)

	// Return
	if coreErr != nil {
		createReplRuntimeErrorResponse(c, coreErr)
		return
	}
	createReplRuntimeResponse(c, *coreResp)
}

// Ensure craft user input to service core standard
func createReplRuntimeRequest(c *gin.Context) (*core.CreateReplRuntimeRequest, error) {
	// Input validate
	var userInput CreateReplRuntimeRequest
	err := c.BindJSON(&userInput)
	if err != nil {
		fmt.Println("[ERROR] API error, can't parse JSON value, got: ", err.Error())
		return nil, fmt.Errorf("Can't parse JSON value")
	}

	req := &core.CreateReplRuntimeRequest{}
	return req, nil
}

// Ensure output to api endpoit return to match API schema standard
func createReplRuntimeResponse(c *gin.Context, coreResp core.CreateReplRuntimeResponseSuccess) {
	// Return
	resp := CreateReplRuntimeResponseSuccess{
		RuntimeId: coreResp.RuntimeId,
	}
	c.JSON(http.StatusOK, resp)
}

// Ensure output to api endpoit return to match API schema standard
func createReplRuntimeErrorResponse(c *gin.Context, error common.ErrorResponse) {
	switch error := error.(type) {
	default:
		resp := EvaluateResponseError{
			Type:    error.GetType(),
			Code:    error.GetCode(),
			Message: error.GetMessage(),
		}
		c.JSON(resp.Type, resp)
	}
}
