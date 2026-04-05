package repl

import (
	"fmt"
	"interingo/pkg/evaluator"
	"interingo/pkg/object"
	"io"
)

type Print struct {
	env *object.Environment
	r   *Repl
}

func (b *Print) Description() string { return "Print into stdout" }
func (b *Print) Func(env *object.Environment) object.Object {
	value, ok := env.Get(evaluator.REST_ARGS)
	if !ok {
		return &object.Error{Message: "Print call doesn't correctly provide rest args"}
	}

	arrValues, ok := value.(*object.Array)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("Print call doesn't correctly provide rest args, got: %v, want: %v", value.Type(), object.ARRAY_OBJ),
		}
	}

	for i, v := range arrValues.Value {
		if i > 0 {
			io.WriteString(b.r.out, " ")
		}
		fmt.Fprintf(b.r.out, "%s", v.Inspect())
	}
	io.WriteString(b.r.out, "\n")

	return nil
}

func (b *Print) Parameters() object.Parameters {
	return object.Parameters{
		Rest: true,
	}
}
func (b *Print) Env() *object.Environment { return b.env }
func (b *Print) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *Print) Inspect() string          { return object.BuiltInInspect(b) }
