package runtime

import (
	"fmt"
	"interingo/pkg/evaluator"
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

func (b *ToggleVerbose) Parameters() object.Parameters {
	return object.Parameters{}
}
func (b *ToggleVerbose) Env() *object.Environment { return nil }
func (b *ToggleVerbose) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *ToggleVerbose) Inspect() string          { return object.BuiltInInspect(b) }

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
func (b *GetRuntimeInfo) Inspect() string          { return object.BuiltInInspect(b) }

// Rely on runtime print for correctly display the help
type Help struct {
	Core *Core
	env  *object.Environment
}

func (b *Help) Description() string {
	return "Print all builtin command"
}

func (b *Help) Func(env *object.Environment) object.Object {
	printBuiltin, ok := env.Get("print")
	if !ok {
		return &object.Error{
			Message: "Help rely's `print` builtin doesn't existed",
		}
	}
	bi, ok := printBuiltin.(object.BuiltIn)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Help rely's `print` builtin, but got %v print instead", printBuiltin.Type()),
		}
	}
	// Construct Builtin - Helper
	infos := b.Core.Help()
	// Print the value
	printEnv := object.NewEnclosedEnvironment(env)
	printEnv.Set(evaluator.REST_ARGS, &object.Array{Value: []object.Object{&object.String{Value: "Built-in commands:\n" + infos}}})
	bi.Func(printEnv)
	return nil
}

func (b *GetRuntimeInfo) WildcastParameters() bool {
	return false
}

func (b *Help) Parameters() object.Parameters {
	return object.Parameters{}
}

func (b *Help) Env() *object.Environment { return b.env }
func (b *Help) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *Help) Inspect() string          { return object.BuiltInInspect(b) }
