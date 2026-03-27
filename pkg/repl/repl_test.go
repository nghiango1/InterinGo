package repl

import(
	"interingo/pkg/service/core"
	"io"
	"testing"
)

func TestRepl_printVerboseInfomation(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		in  io.Reader
		out io.Writer
		// Named input parameters for target function.
		info *core.VerboseInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepl(tt.in, tt.out)
			r.printVerboseInfomation(tt.info)
		})
	}
}

