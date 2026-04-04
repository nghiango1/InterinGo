package runtime

import (
	"interingo/pkg/object"
)

// Core Runtime specific builtin
type ToggleVerbose struct {
	Core *Core
}

func (b *ToggleVerbose) Description() string {
	return "Toggle more verbose information to be return by the core runtime"
}
func (b *ToggleVerbose) Func(env *object.Environment) object.Object {
	b.Core.ToggleVerbose()
	return &object.Boolean{Value: b.Core.Verbose}
}

func (b *ToggleVerbose) WildcastParameters() bool {
	return false
}
func (b *ToggleVerbose) Parameters() object.Parameters {
	return object.Parameters{}
}
func (b *ToggleVerbose) Env() *object.Environment { return nil }
func (b *ToggleVerbose) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *ToggleVerbose) Inspect() string          { return "BuiltIn: " + b.Description() }

type GetRuntimeInfo struct {
	Core *Core
}

func (b *GetRuntimeInfo) Description() string {
	return "Get current core runtime informations (stdout internal only)"
}

func (b *GetRuntimeInfo) Func(env *object.Environment) object.Object {
	return &object.Boolean{Value: true}
}

func (b *GetRuntimeInfo) Parameters() object.Parameters {
	return object.Parameters{}
}
func (b *GetRuntimeInfo) Env() *object.Environment { return nil }
func (b *GetRuntimeInfo) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *GetRuntimeInfo) Inspect() string          { return "BuiltIn: " + b.Description() }
