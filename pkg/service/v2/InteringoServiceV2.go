package v2

import (
	"interingo/pkg/service/core"

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

