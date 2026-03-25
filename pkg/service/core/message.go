package core

type EvaluateRequest struct {
	Data string
}

type EvaluateResponseSuccess struct {
	Output  string
	Verbose *VerboseInfo
}

type ParserErrorResponse struct {
	Errors  []string
	Type    string
	Code    string
	Message string
}

func (resp *ParserErrorResponse) Error()

type EvaluateResponse interface {
	Success() EvaluateResponse
	Error() ErrorResponse
}

type ErrorResponse interface {
	Type() string
	Code() string
	Message() string
}

type VerboseInfo map[string]any
