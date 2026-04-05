package repl_test

import(
	"interingo/pkg/object"
	"interingo/pkg/repl"
	"testing"
)

func TestPrint_Builtin(t *testing.T) {
	var _ object.BuiltIn = &repl.Print{}
}

