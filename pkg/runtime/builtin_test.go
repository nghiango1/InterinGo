package runtime

import (
	"interingo/pkg/object"
	"interingo/pkg/test"
	"testing"
)

func TestBuiltin(t *testing.T) {
	c := NewCore(EMBED, nil)
	var toggleVerbose object.BuiltIn = &ToggleVerbose{Core: c}

	prior := c.Verbose
	toggleVerbose.Func(nil)
	if c.Verbose == prior {
		t.Fatalf("Verbose doesn't match")
	}
	env := object.NewEnvironment()
	env.Set("print", &test.MockBuiltinImpl{
		MockParameters: object.Parameters{},
		MockEnv:        env,
	})
	var getRuntimeInfo object.BuiltIn = &GetRuntimeInfo{Core: c, env: env}
	// return nil, as long as there no error here then it is fine
	getRuntimeInfo.Func(env)
}
