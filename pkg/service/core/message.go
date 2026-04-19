package core

import (
	"interingo/pkg/parser"
	"interingo/pkg/runtime"
	"interingo/pkg/service/common"
	"strings"
)

type CreateReplRuntimeRequest struct{}

type CreateReplRuntimeResponseSuccess struct {
	RuntimeId string `json:"runtimeId"`
}

type EvaluateRequest struct {
	RuntimeId string `json:"id"`
	Data      string `json:"data"`
}

type EvaluateResponse struct {
	Success *EvaluateResponseSuccess
	Error   common.ErrorResponseInterface
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

// Websocket
type WebsocketConnectSuccess struct {
	Type   WebsocketMessage `json:"type"` // type already a keyword
	ConnId string           `json:"connId"`
}

func NewWebsocketConnectSuccess(connId string) *WebsocketConnectSuccess {
	return &WebsocketConnectSuccess{
		Type:   WS_OPEN,
		ConnId: connId,
	}
}

type WebsocketConnectError struct {
	Type  WebsocketMessage `json:"type"` // type already a keyword
	Error string           `json:"error"`
}

func NewWebsocketConnectError(error string) *WebsocketConnectError {
	return &WebsocketConnectError{
		Type:  WS_ERROR,
		Error: error,
	}
}

type PrintMessageEventData struct {
	Type    WebsocketMessage `json:"type"` // type already a keyword
	Message string           `json:"message"`
}

func NewPrintMessageEventData(message string) *PrintMessageEventData {
	return &PrintMessageEventData{
		Type:    WS_PRINT,
		Message: message,
	}
}

type WebsocketMessage string

const (
	WS_UNKNOW = WebsocketMessage("ws_unknow")
	WS_OPEN   = WebsocketMessage("ws_open")
	WS_ERROR  = WebsocketMessage("ws_error")
	WS_PRINT  = WebsocketMessage("ws_print")
)
