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
	Error   common.ErrorResponse
	Verbose *runtime.VerboseInfo `json:"verbose"`
}

type EvaluateResponseSuccess struct {
	Output *string `json:"output,omitempty"`
}

// Share common interface of common.ErrorResponseInterface
type ParserErrorResponse struct {
	Message string               `json:"message"`
	Errors  []parser.ParserError `json:"error,omitempty"`
}

func NewParserErrorResponse(message string, errors []parser.ParserError) *ParserErrorResponse {
	// Make sure Parser errors is clear at start, this help with UI display
	// Could have better way, this work for now
	if message == "" {
		message = "PARSER ERROR: provided code was invalid,"
	} else if !strings.HasPrefix(message, "PARSER ERROR") {
		message = "PARSER ERROR: " + message
	}

	return &ParserErrorResponse{
		Message: message,
		Errors:  errors,
	}
}

func (e *ParserErrorResponse) GetType() int       { return 400 }
func (e *ParserErrorResponse) GetCode() string    { return "parser_error" }
func (e *ParserErrorResponse) GetMessage() string { return e.Message }

// Share common interface of common.ErrorResponseInterface
type EvalErrorResponse struct {
	Message string `json:"message"`
}

func NewEvalErrorResponse(message string) *EvalErrorResponse {
	// Make sure Parser errors is clear at start, this help with UI display
	// Could have better way, this work for now
	if message == "" {
		message = "ERROR: provided code was invalid,"
	} else if !strings.HasPrefix(message, "ERROR") {
		message = "ERROR: " + message
	}

	return &EvalErrorResponse{
		Message: message,
	}
}

func (e *EvalErrorResponse) GetType() int       { return 400 }
func (e *EvalErrorResponse) GetCode() string    { return "eval_error" }
func (e *EvalErrorResponse) GetMessage() string { return e.Message }

// Websocket
type WebsocketConnectSuccess struct {
	Type   WebsocketMessageResponseType `json:"type"` // type already a keyword
	ConnId string                       `json:"connId"`
}

func NewWebsocketConnectSuccess(connId string) *WebsocketConnectSuccess {
	return &WebsocketConnectSuccess{
		Type:   WS_OPEN,
		ConnId: connId,
	}
}

type WebsocketConnectError struct {
	Type  WebsocketMessageResponseType `json:"type"` // type already a keyword
	Error string                       `json:"error"`
}

func NewWebsocketConnectError(error string) *WebsocketConnectError {
	return &WebsocketConnectError{
		Type:  WS_ERROR,
		Error: error,
	}
}

type PrintMessageEventData struct {
	Type    WebsocketMessageResponseType `json:"type"` // type already a keyword
	Message string                       `json:"message"`
}

func NewPrintMessageEventData(message string) *PrintMessageEventData {
	return &PrintMessageEventData{
		Type:    WS_PRINT,
		Message: message,
	}
}

type WebsocketMessageResponseType string

const (
	WS_UNKNOW = WebsocketMessageResponseType("ws_unknow")
	WS_OPEN   = WebsocketMessageResponseType("ws_open")
	WS_ERROR  = WebsocketMessageResponseType("ws_error")
	WS_PRINT  = WebsocketMessageResponseType("ws_print")
)

type WebsocketRequest interface {
	Type() WebsocketRequestType
}

type ReplBindRequest struct {
	MessageType WebsocketRequestType `json:"type"`
	RuntimeId   string               `json:"runtimeId"`
}

func (m *ReplBindRequest) Type() WebsocketRequestType {
	return m.MessageType
}

type WebsocketRequestType string

const (
	REPL_BIND = WebsocketRequestType("repl_bind")
)
