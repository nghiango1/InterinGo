package core

import (
	"log"

	"interingo/pkg/runtime"
	"interingo/pkg/service/common"

	"github.com/google/uuid"
)

func (c *ServiceCore) CreateReplRuntime(req CreateReplRuntimeRequest) (*CreateReplRuntimeResponseSuccess, common.ErrorResponseInterface) {
	c.muConnClients.Lock()
	defer c.muConnClients.Unlock()

	runtimeId := uuid.New().String()
	_, ok := c.runtimeCores[runtimeId]
	if ok {
		log.Printf("[ERROR] ConnId collision, should not be possible")
		return nil, common.NewErrorResponse(500)
	}

	evalCore := runtime.NewCore(runtime.EMBED, nil)

	c.runtimeCores[runtimeId] = &ReplRuntime{
		id:   runtimeId,
		core: evalCore,
	}

	// If both err and res from Eval is nil, there some thing wrong
	return &CreateReplRuntimeResponseSuccess{
		RuntimeId: runtimeId,
	}, nil
}
