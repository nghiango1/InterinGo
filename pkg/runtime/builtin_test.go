package runtime_test

import (
	"interingo/pkg/object"
	"interingo/pkg/runtime"
	"testing"
)

func TestBuiltin(t *testing.T) {
	c := runtime.NewCore(runtime.EMBED, nil)
	var toggleVerbose object.BuiltIn = &runtime.ToggleVerbose{Core: c}

	prior := c.Verbose
	toggleVerbose.Func(nil)
	if c.Verbose == prior {
		t.Fatalf("Verbose doesn't match")
	}
	var getRuntimeInfo object.BuiltIn = &runtime.GetRuntimeInfo{Core: c}
	got := getRuntimeInfo.Func(nil)
	if v, ok := got.(*object.Boolean); ok {
		if v.Value != true {
			t.Fatalf("GetRuntimeInfo Return failured vaule. Got %v, want %v", v.Value, true)
		}
	} else {
		t.Fatalf("GetRuntimeInfo Return vaule type not match. Got %v, want %v", got.Type(), object.BOOLEAN_OBJ)
	}
}
