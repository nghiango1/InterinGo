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
func (b *SystemExit) Func(env *object.Environment) object.Object {
	code, ok := env.Get("code")
	if !ok {
		os.Exit(1)
	}

	val, ok := code.(*object.Integer)
	if !ok {
		os.Exit(1)
	}

	os.Exit(int(val.Value))
	return &object.Null{}
}

func (b *SystemExit) Parameters() object.Parameters {
	return object.Parameters{
		Default: []object.DefaultParameter{
			{
				Key:   &ast.Identifier{Value: "code"},
				Value: &object.Integer{Value: 0},
			},
		},
	}
}
func (b *SystemExit) Env() *object.Environment { return b.Environment }
func (b *SystemExit) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *SystemExit) Inspect() string          { return object.BuiltInInspect(b) }
