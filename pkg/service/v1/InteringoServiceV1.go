package v1

import (
	"fmt"
	"interingo/pkg/service/common"
	"interingo/pkg/service/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InteringoServiceV1 interface {
	EvaluateHandler(c *gin.Context)
}

func NewInteringoServiceV1(core *core.ServiceCore) InteringoServiceV1 {
	return &interingoServiceV1Impl{
		serviceCore: core,
	}
}

type interingoServiceV1Impl struct {
	serviceCore *core.ServiceCore
}

// API v1 api/evalutate
func (s *interingoServiceV1Impl) EvaluateHandler(c *gin.Context) {
	// Input validate
	var req core.EvaluateRequest
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

	res, evalErr := s.serviceCore.EvaluateHandler(req)

	// Return
	if evalErr != nil {
		c.JSON(evalErr.GetType(), evalErr)
	} else if res != nil {
		c.JSON(http.StatusOK, res)
	}
}
