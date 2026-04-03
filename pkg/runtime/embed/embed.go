package embed

import (
	"interingo/pkg/ast"
	"interingo/pkg/object"
)

type SystemExit struct {
	Environment *object.Environment
}

func (b *SystemExit) Description() string { return "Exit the program" }

// Embed return specific SystemExit object, which leave the Eval runtime
// to correctly route to the actual Exit handler
func (b *SystemExit) Func(env *object.Environment) object.Object {
	code, ok := env.Get("code")
	if !ok {
		return &object.SystemExit{Code: 1}
	}

	val, ok := code.(*object.Integer)
	if !ok {
		return &object.SystemExit{Code: 1}
	}

	return &object.SystemExit{Code: int(val.Value)}
}
func (b *SystemExit) Parameters() []*ast.Identifier {
	return []*ast.Identifier{
		{
			Value: "code",
		},
	}
}
func (b *SystemExit) Env() *object.Environment { return b.Environment }
func (b *SystemExit) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *SystemExit) Inspect() string          { return "BuiltIn: " + b.Description() }
