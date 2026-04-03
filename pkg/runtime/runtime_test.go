package runtime_test

import (
	"fmt"
	"interingo/pkg/runtime"
	"testing"
)

func newString(s any) *string {
	out := fmt.Sprintf("%v", s)
	return &out
}

func evalResponseSuccessCompare(got *runtime.EvalResponseSuccess, want *runtime.EvalResponseSuccess) bool {
	if want == nil {
		return got == nil
	}

	// want is not nil, got is nil
	if got == nil {
		return false
	}

	if want.Output == nil {
		return got.Output == nil
	}

	// want.Output is not nil, got.Output is nil
	if got.Output == nil {
		return false
	}

	return *(want.Output) == *(got.Output)
}

func TestCore_Eval(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req   runtime.EvalRequest
		want  *runtime.EvalResponseSuccess
		want2 *runtime.EvalResponseError
		want3 *runtime.VerboseInfo
	}{
		{
			"LetStatements",
			runtime.EvalRequest{Data: "let identity = fn(x) { return x; }; identity(identity(5));"},
			&runtime.EvalResponseSuccess{Output: newString(5)},
			nil,
			nil,
		},
		{
			"LetStatements",
			runtime.EvalRequest{Data: "let identity = fn(x) { x; }; identity(5);"},
			&runtime.EvalResponseSuccess{Output: newString(5)},
			nil,
			nil,
		},
		{
			"FunctionApplication",
			runtime.EvalRequest{Data: "fn(x) { x; }(5)"},
			&runtime.EvalResponseSuccess{Output: newString(5)},
			nil,
			nil,
		},
		{
			"FunctionApplication",
			runtime.EvalRequest{Data: "let double = fn(x) { x * 2; }; double(5);"},
			&runtime.EvalResponseSuccess{Output: newString(10)},
			nil,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := runtime.NewCore(runtime.NATIVE, nil)
			got, _, _ := c.Eval(tt.req)
			if !evalResponseSuccessCompare(got, tt.want) {
				t.Errorf("Eval() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}

func TestCore_EvalForceVerbose(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req  runtime.EvalRequest
		want *runtime.EvalResponseSuccess
	}{
		{
			"Verbose_LetStatements",
			runtime.EvalRequest{Data: "let identity = fn(x) { return x; }; identity(identity(5));"},
			&runtime.EvalResponseSuccess{Output: newString(5)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := runtime.NewCore(runtime.EMBED, nil)
			c.ToggleVerbose()
			got, _, ver := c.Eval(tt.req)
			if !evalResponseSuccessCompare(got, tt.want) {
				t.Errorf("Eval() = %v, want %v", got.String(), tt.want.String())
			}
			if ver == nil {
				t.Errorf("Verbose is nil")
			}
		})
	}
}

func TestCore_EvalBuiltinVerbose(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req  runtime.EvalRequest
		want *runtime.EvalResponseSuccess
	}{
		{
			"Verbose_LetStatements",
			runtime.EvalRequest{Data: "toggleVerbose(); let identity = fn(x) { return x; }; identity(identity(5));"},
			&runtime.EvalResponseSuccess{Output: newString(5)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := runtime.NewCore(runtime.EMBED, nil)
			got, _, ver := c.Eval(tt.req)
			if !evalResponseSuccessCompare(got, tt.want) {
				t.Errorf("Eval() = %v, want %v", got.String(), tt.want.String())
			}
			if ver == nil {
				t.Errorf("Verbose is nil")
			}
		})
	}
}


func TestCore_Eval_ParserErr(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req runtime.EvalRequest
	}{
		{
			"ParserError_LetStatements",
			runtime.EvalRequest{Data: "let identity = fn(x) { return ; identity(identity(5));"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := runtime.NewCore(runtime.EMBED, nil)
			c.ToggleVerbose()
			_, err, _ := c.Eval(tt.req)
			if err == nil {
				t.Errorf("Error is nil")
			}
		})
	}
}

func TestCore_EvalSystemExit(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req      runtime.EvalRequest
		res      *runtime.EvalResponseSuccess
		exitCode int
		state    runtime.RuntimeState
	}{
		{
			"ExitStatements",
			runtime.EvalRequest{Data: "let identity = fn(x) { return x; }; exit(2); identity(identity(5));"},
			nil,
			2,
			runtime.RUNTIME_END,
		},
		{
			"ExitStatements 2",
			runtime.EvalRequest{Data: "let identity = fn(x) { return x; }; exit(identity); identity(identity(5));"},
			nil,
			1,
			runtime.RUNTIME_END,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := runtime.NewCore(runtime.EMBED, nil)
			c.ToggleVerbose()
			got, err, _ := c.Eval(tt.req)
			if !evalResponseSuccessCompare(got, tt.res) {
				t.Errorf("Eval() = %v, want %v", got.String(), tt.res.String())
			}
			if err == nil {
				t.Errorf("Error is nil")
			}
			if err.SystemExit == nil {
				t.Errorf("Error return with SYSTEM_EXIT is nil")
			}
			if err.SystemExit.Code != tt.exitCode {
				t.Errorf("Error return with SYSTEM_EXIT but code doesn't match. Want %d got %d", tt.exitCode, err.SystemExit.Code)
			}
			if c.State() != tt.state {
				t.Errorf("Runtime state doesn't match. Want %d got %d", tt.state, c.State())
			}

			_, err, _ = c.Eval(runtime.EvalRequest{Data: "1"})
			if err == nil {
				t.Errorf("Expect error when Runtime Eval after called Exit()")
			}
			if err.Error == nil {
				t.Errorf("Expect error when Runtime Eval after called Exit()")
			}
		})
	}
}
