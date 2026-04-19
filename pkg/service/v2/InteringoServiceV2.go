package v2

import (
	"fmt"
	"interingo/pkg/service/common"
	"interingo/pkg/service/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InteringoServiceV2 interface {
	EvaluateHandler(c *gin.Context)
	CreateReplRuntimeHandler(c *gin.Context)
}

func NewInteringoServiceV2(core *core.ServiceCore) InteringoServiceV2 {
	return &interingoServiceV2Impl{
		serviceCore: core,
	}
}

type interingoServiceV2Impl struct {
	serviceCore *core.ServiceCore
}

// API v1 api/repl/:id/evalutate
func (s *interingoServiceV2Impl) EvaluateHandler(c *gin.Context) {
	runtimeId := c.Param("id")
	// Input validate
	var req core.EvaluateRequest
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println("[ERRPR] API error, can't parse JSON value, got: ", err.Error())
		errorResp := common.NewBadRequestErrorResponse("Invalid JSON", nil)
		c.JSON(http.StatusBadRequest, errorResp)
	}

	req.RuntimeId = runtimeId

	if s.serviceCore == nil {
		fmt.Println("[ERRPR] API error, serviceCore didn't init yet")
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500))
	}

	res, evalErr := s.serviceCore.EvaluateHandlerV2(req)

	// Return
	if evalErr != nil {
		c.JSON(evalErr.GetType(), evalErr)
	} else if res != nil {
		c.JSON(http.StatusOK, res)
	}
}

// API v1 api/repl/:id/evalutate
func (s *interingoServiceV2Impl) CreateReplRuntimeHandler(c *gin.Context) {
	// Input validate
	var req core.CreateReplRuntimeRequest
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println("[ERRPR] API error, can't parse JSON value, got: ", err.Error())
		errorResp := common.NewBadRequestErrorResponse("Invalid JSON", nil)
		c.JSON(http.StatusBadRequest, errorResp)
	}

	if s.serviceCore == nil {
		fmt.Println("[ERRPR] API error, serviceCore didn't init yet")
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500))
	}

	res, evalErr := s.serviceCore.CreateReplRuntime(req)

	// Return
	if evalErr != nil {
		c.JSON(evalErr.GetType(), evalErr)
	} else if res != nil {
		c.JSON(http.StatusOK, res)
	}
}
