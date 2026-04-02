package object

import (
	"bytes"
	"fmt"
	"interingo/pkg/ast"
	"os"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILT_IN_OBJ     = "BUILT_IN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type BuiltIn interface {
	Object
	Description() string
	Func() func(env *Environment) Object
	Parameters() []*ast.Identifier
	Env() *Environment
}

type SystemExit struct {
	Environment *Environment
}

func (b *SystemExit) Description() string { return "Exit the program" }
func (b *SystemExit) Func() func(env *Environment) Object {
	return func(env *Environment) Object {
		code, ok := env.Get("code")
		if !ok {
			os.Exit(0)
		}

		val, ok := code.(*Integer)
		if !ok {
			os.Exit(0)
		}

		os.Exit(int(val.Value))
		return &Null{}
	}
}
func (b *SystemExit) Parameters() []*ast.Identifier {
	return []*ast.Identifier{
		{
			Value: "code",
		},
	}
}
func (b *SystemExit) Env() *Environment { return b.Environment }
func (b *SystemExit) Type() ObjectType  { return BUILT_IN_OBJ }
func (b *SystemExit) Inspect() string   { return "BuiltIn: " + b.Description() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return fmt.Sprintf("%s", rv.Value.Inspect()) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())
	return out.String()
}
