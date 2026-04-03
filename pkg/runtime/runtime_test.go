package runtime

import (
	"fmt"
	"testing"
)

func newString(s any) *string {
	out := fmt.Sprintf("%v", s)
	return &out
}

func evalResponseSuccessCompare(got *EvalResponseSuccess, want *EvalResponseSuccess) bool {
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
		req   EvalRequest
		want  *EvalResponseSuccess
		want2 *EvalResponseError
		want3 *VerboseInfo
	}{
		{
			"LetStatements",
			EvalRequest{Data: "let identity = fn(x) { return x; }; identity(identity(5));"},
			&EvalResponseSuccess{Output: newString(5)},
			nil,
			nil,
		},
		{
			"LetStatements",
			EvalRequest{Data: "let identity = fn(x) { x; }; identity(5);"},
			&EvalResponseSuccess{Output: newString(5)},
			nil,
			nil,
		},
		{
			"FunctionApplication",
			EvalRequest{Data: "fn(x) { x; }(5)"},
			&EvalResponseSuccess{Output: newString(5)},
			nil,
			nil,
		},
		{
			"FunctionApplication",
			EvalRequest{Data: "let double = fn(x) { x * 2; }; double(5);"},
			&EvalResponseSuccess{Output: newString(10)},
			nil,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCore(NATIVE)
			got, _, _ := c.Eval(tt.req)
			if !evalResponseSuccessCompare(got, tt.want) {
				t.Errorf("Eval() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}

func TestCore_Eval_Verbose(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		req  EvalRequest
		want *EvalResponseSuccess
	}{
		{
			"Verbose_LetStatements",
			EvalRequest{Data: "let identity = fn(x) { return x; }; identity(identity(5));"},
			&EvalResponseSuccess{Output: newString(5)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCore(EMBED)
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
