package core

import (
	"fmt"
	"io"
	"log"
	"strings"

	"interingo/pkg/evaluator"
	"interingo/pkg/object"
)

type PrintMessageEventData struct {
	Message string `json:"message"`
}

func (c *ServiceCore) Print(mes string) {
	if c.conn == nil {
		log.Printf("ERROR: Try to print without any client connection")
		return
	}
	c.conn.WriteJSON(PrintMessageEventData{
		Message: mes,
	})
}

type PrintBuiltin struct {
	env  *object.Environment
	core *ServiceCore
}

func (b *PrintBuiltin) Description() string { return "PrintBuiltin into stdout" }
func (b *PrintBuiltin) Func(env *object.Environment) object.Object {
	value, ok := env.Get(evaluator.REST_ARGS)
	if !ok {
		return &object.Error{Message: "PrintBuiltin call doesn't correctly provide rest args"}
	}

	arrValues, ok := value.(*object.Array)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("PrintBuiltin call doesn't correctly provide rest args, got: %v, want: %v", value.Type(), object.ARRAY_OBJ),
		}
	}

	message := &strings.Builder{}
	for i, v := range arrValues.Value {
		if i > 0 {
			io.WriteString(message, " ")
		}
		fmt.Fprintf(message, "%s", v.Inspect())
	}
	io.WriteString(message, "\n")

	b.core.Print(message.String())
	return nil
}

func (b *PrintBuiltin) Parameters() object.Parameters {
	return object.Parameters{
		Rest: true,
	}
}
func (b *PrintBuiltin) Env() *object.Environment { return b.env }
func (b *PrintBuiltin) Type() object.ObjectType  { return object.BUILT_IN_OBJ }
func (b *PrintBuiltin) Inspect() string          { return object.BuiltInInspect(b) }
