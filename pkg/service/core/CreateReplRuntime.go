package core

import (
	"log"
	"time"

	"interingo/pkg/runtime"
	"interingo/pkg/service/common"

	"github.com/google/uuid"
)

const (
	alivePeriod = 60 * time.Second
)

func (c *ServiceCore) CreateReplRuntime(req CreateReplRuntimeRequest) (*CreateReplRuntimeResponseSuccess, common.ErrorResponse) {
	c.muConnClients.Lock()
	defer c.muConnClients.Unlock()

	runtimeId := uuid.New().String()
	_, ok := c.runtimeCores[runtimeId]
	if ok {
		log.Printf("[ERROR] ConnId collision, should not be possible")
		return nil, common.NewErrorResponse(500)
	}

	evalCore := runtime.NewCore(runtime.EMBED, nil)

	runtime := &ReplRuntime{
		id:   runtimeId,
		core: evalCore,
	}
	c.runtimeCores[runtimeId] = runtime

	// With out Client connection, we need to clean up base from the request
	// for evaluate of this REPL session
	// Start a goroutine to send pings.
	go func() {
		// good enough, timer check every 60 second
		ticker := time.NewTicker(alivePeriod)
		defer ticker.Stop()
		for range ticker.C {
			if runtime.connId != "" {
				// conn client will help with clean up
				break
			}
			// else check if lastUpdate + alivePeriod < current time.
			// If so clean up is perform
			if time.Now().After(runtime.lastUpdate.Add(alivePeriod)) {
				c.muConnClients.Lock()
				defer c.muConnClients.Unlock()
				delete(c.runtimeCores, runtime.id)
			}
		}
	}()

	// If both err and res from Eval is nil, there some thing wrong
	return &CreateReplRuntimeResponseSuccess{
		RuntimeId: runtimeId,
	}, nil
}
