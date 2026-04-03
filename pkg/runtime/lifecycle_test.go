package runtime_test

import (
	"interingo/pkg/runtime"
	"testing"
)

type defaultLifeCycleHandlerImpt struct {
	ExitCalled int
}

func (lc *defaultLifeCycleHandlerImpt) Exit() {
	lc.ExitCalled += 1
}

func TestCore_LifeCycleExit(t *testing.T) {
	lc := &defaultLifeCycleHandlerImpt{}
	c := runtime.NewCore(runtime.EMBED, lc)
	c.ToggleVerbose()
	res, err, _ := c.Eval(runtime.EvalRequest{Data: "exit(1);"})
	if err == nil {
		t.Errorf("Exit not working for EMBED, got res: %v, err: %v", res, err)
	}
	if lc.ExitCalled != 1 {
		t.Errorf("LifeCycle exit handler not called")
	}
}
