// Should be generate with mock test. Eg: Mockery later
package test

import (
	"interingo/pkg/evaluator"
	"interingo/pkg/object"
)

type MockBuiltinImpl struct {
	MockParameters object.Parameters
	MockEnv        *object.Environment
}

func (b *MockBuiltinImpl) Description() string      { return "Mock implement" }
func (b *MockBuiltinImpl) Env() *object.Environment { return b.MockEnv }
func (b *MockBuiltinImpl) Func(env *object.Environment) object.Object {
	b.MockEnv = env
	return evaluator.NULL
}
func (b *MockBuiltinImpl) Inspect() string               { return "BUILT_IN: " + b.Description() }
func (b *MockBuiltinImpl) Type() object.ObjectType       { return object.BUILT_IN_OBJ }
func (b *MockBuiltinImpl) Parameters() object.Parameters { return b.MockParameters }
