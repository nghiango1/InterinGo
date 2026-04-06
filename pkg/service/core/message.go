package core

import (
	"interingo/pkg/parser"
	"interingo/pkg/runtime"
)

type EvaluateRequest struct {
	Data string `json:"data"`
}

type EvaluateResponseSuccess struct {
	Output  *string              `json:"output,omitempty"`
	Verbose *runtime.VerboseInfo `json:"verbose,omitempty"`
}

// Share common interface of common.ErrorResponseInterface
type ParserErrorResponse struct {
	Type    int                  `json:"type"` // type already a keyword
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Errors  []parser.ParserError `json:"error,omitempty"`
	Verbose *runtime.VerboseInfo `json:"verbose,omitempty"`
}

func NewParserErrorResponse(message string, errors []parser.ParserError, verbose *runtime.VerboseInfo) *ParserErrorResponse {
	if message == "" {
		message = "Parse error: provided code was invalid,"
	}

	return &ParserErrorResponse{
		Type:    400,
		Code:    "parser_error",
		Message: message,
		Errors:  errors,
		Verbose: verbose,
	}
}

func (e *ParserErrorResponse) GetType() int       { return e.Type }
func (e *ParserErrorResponse) GetCode() string    { return e.Code }
func (e *ParserErrorResponse) GetMessage() string { return e.Message }
