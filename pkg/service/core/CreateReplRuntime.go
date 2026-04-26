package core

import (
	"log/slog"

	"interingo/pkg/runtime"
	"interingo/pkg/service/common"

	"github.com/google/uuid"
)

func (c *ServiceCore) CreateReplRuntime(req CreateReplRuntimeRequest) (*CreateReplRuntimeResponseSuccess, common.ErrorResponse) {
	slog.Debug("CreateReplRuntime request lock muConnClients")
	c.muConnClients.Lock()
	slog.Debug("CreateReplRuntime take lock muConnClients")
	defer c.muConnClients.Unlock()
	defer slog.Debug("CreateReplRuntime release lock muConnClients")

	runtimeId := uuid.New().String()
	_, ok := c.runtimeCores[runtimeId]
	if ok {
		slog.Error("ConnId collision, should not be possible")
		return nil, common.NewErrorResponse(500)
	}

	evalCore := runtime.NewCore(runtime.EMBED, nil)

	runtime := &ReplRuntime{
		id:   runtimeId,
		core: evalCore,
	}
	c.runtimeCores[runtimeId] = runtime

	// If both err and res from Eval is nil, there some thing wrong
	return &CreateReplRuntimeResponseSuccess{
		RuntimeId: runtimeId,
	}, nil
}
