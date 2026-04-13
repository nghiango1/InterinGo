package core

import (
	"interingo/pkg/parser"
	"interingo/pkg/runtime"
	"strings"
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
	// Make sure Parser errors is clear at start, this help with UI display
	// Could have better way, this work for now
	if message == "" {
		message = "PARSER ERROR: provided code was invalid,"
	} else if !strings.HasPrefix(message, "PARSER ERROR") {
		message = "PARSER ERROR: " + message
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

// Share common interface of common.ErrorResponseInterface
type EvalErrorResponse struct {
	Type    int                  `json:"type"` // type already a keyword
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Verbose *runtime.VerboseInfo `json:"verbose"`
}

func NewEvalErrorResponse(message string, verbose *runtime.VerboseInfo) *EvalErrorResponse {
	// Make sure Parser errors is clear at start, this help with UI display
	// Could have better way, this work for now
	if message == "" {
		message = "ERROR: provided code was invalid,"
	} else if !strings.HasPrefix(message, "ERROR") {
		message = "ERROR: " + message
	}

	return &EvalErrorResponse{
		Type:    400,
		Code:    "eval_error",
		Message: message,
	}
}

func (e *EvalErrorResponse) GetType() int       { return e.Type }
func (e *EvalErrorResponse) GetCode() string    { return e.Code }
func (e *EvalErrorResponse) GetMessage() string { return e.Message }
