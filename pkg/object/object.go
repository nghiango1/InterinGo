package object

import (
	"bytes"
	"fmt"
	"interingo/pkg/ast"
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
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	BUILT_IN_OBJ     = "BUILT_IN"
	SYSTEM_EXIT_OBJ  = "SYSTEM_EXIT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Runtime Builtin, object to inject extra runtime functionality
// Should be implement and added when create a Runtime
// - Eg: Exit()
type BuiltIn interface {
	Object
	Description() string
	Func(env *Environment) Object
	Parameters() Parameters
	Env() *Environment
}

func FnParamsInspect(fp Parameters) string {
	params := []string{}
	for _, p := range fp.Standard {
		params = append(params, p.String())
	}
	for _, p := range fp.Default {
		params = append(params, p.String())
	}
	if fp.Rest {
		params = append(params, "*args")
	}

	return strings.Join(params, ", ")
}

func BuiltInInspect(f BuiltIn) string {
	var out bytes.Buffer
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(FnParamsInspect(f.Parameters()))
	out.WriteString(") : ")
	out.WriteString(f.Description())
	return out.String()
}

type Parameters struct {
	Standard []*ast.Identifier
	Default  []DefaultParameter
	Rest     bool
}

type DefaultParameter struct {
	Key   *ast.Identifier
	Value Object
}

func (d *DefaultParameter) String() string {
	return fmt.Sprintf("%v=%v", d.Key.String(), d.Value.Inspect())
}

// To stop the runtime, should atc like Error but have higher prioirty
type SystemExit struct {
	Code int
}

func (e *SystemExit) Type() ObjectType { return ERROR_OBJ }
func (e *SystemExit) Inspect() string  { return fmt.Sprintf("SYSTEM_EXIT: %d", e.Code) }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Array struct {
	Value []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var res strings.Builder
	res.WriteString("[")
	for i, obj := range a.Value {
		if i > 0 {
			res.WriteString(", ")
		}
		res.WriteString(obj.Inspect())
	}
	res.WriteString("]")
	return res.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string {
	return s.Value
}

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
