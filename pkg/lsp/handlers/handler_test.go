package handlers

import (
	"fmt"
	"interingo/pkg/lexer"
	"interingo/pkg/parser"
	"log/slog"
	"testing"
)

func TestInvalidParsingDianostic(t *testing.T) {
	input := `// 5
let identity = fn(x) {
	// 5
	return x;
	// 5
};
/leta fdaasF a //
return /a as dasd
	let x=/
	 if () \
	 if \ \
identity(5);
`
	l := lexer.New(input)
	p := parser.New(l)
	p.ParseProgram()

	found := getDiagnostic(p.Errors)
	for _, f := range found {
		slog.Debug(fmt.Sprintf("%v", f))
	}
}
