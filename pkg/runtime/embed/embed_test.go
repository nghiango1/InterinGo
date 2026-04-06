package embed_test

import (
	"interingo/pkg/object"
	"interingo/pkg/runtime/embed"
	"testing"
)

func TestSystemExit_Builtin(t *testing.T) {
	var _ object.BuiltIn = &embed.SystemExit{}
}
