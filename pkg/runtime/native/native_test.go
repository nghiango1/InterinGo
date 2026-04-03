package native_test

import(
	"interingo/pkg/object"
	"interingo/pkg/runtime/native"
	"testing"
)

func TestSystemExit_Builtin(t *testing.T) {
	var _ object.BuiltIn = &native.SystemExit{}
}
