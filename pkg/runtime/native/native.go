// Builtin declearation for native session
package native

import (
	"os"

	"interingo/pkg/ast"
	"interingo/pkg/object"
)

type SystemExit struct {
	Environment *object.Environment
}

func (b *SystemExit) Description() string { return "Exit the program" }
func (b *SystemExit) Func() func(env *object.Environment) object.Object {
	return func(env *object.Environment) object.Object {
		code, ok := env.Get("code")
		if !ok {
			os.Exit(0)
		}

		val, ok := code.(*object.Integer)
		if !ok {
			os.Exit(0)
		}

		os.Exit(int(val.Value))
		return &object.Null{}
	}
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
