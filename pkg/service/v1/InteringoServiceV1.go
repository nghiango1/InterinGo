package v1

import (
	"interingo/pkg/service/core"

	"github.com/gin-gonic/gin"
)

// No websocket is supported at this point, so all print and verbose data
// should be within the output response
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

