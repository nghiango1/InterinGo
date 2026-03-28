package core

import "interingo/pkg/token"

type EvaluateRequest struct {
	Data string `json:"data"`
}

type EvaluateResponseSuccess struct {
	Output  *string       `json:"output,omitempty"`
	Verbose *VerboseInfo `json:"verbose,omitempty"`
}

// Share common interface of common.ErrorResponseInterface
type ParserErrorResponse struct {
	Type    string       `json:"type"` // type already a keyword
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Errors  []string     `json:"error"`
	Verbose *VerboseInfo `json:"verbose,omitempty"`
}

func NewParserErrorResponse(message string, errors []string, verbose *VerboseInfo) *ParserErrorResponse {
	if message == "" {
		message = "Parse error: provided code was invalid,"
	}

	return &ParserErrorResponse{
		Type:    "bad_request",
		Code:    "parser_error",
		Message: message,
		Errors:  errors,
		Verbose: verbose,
	}
}

func (e *ParserErrorResponse) GetType() string    { return e.Type }
func (e *ParserErrorResponse) GetCode() string    { return e.Code }
func (e *ParserErrorResponse) GetMessage() string { return e.Message }

type VerboseInfo struct {
	Lexer  LexerInfo  `json:"lexer"`
	Parser ParserInfo `json:"parser"`
}
type LexerInfo struct {
	WhitespaceSkip int                     `json:"whitespace"`
	CommentLine    int                     `json:"comment"`
	Token          map[token.TokenType]int `json:"token"`
}

type ParserInfo struct {
	Ats any `json:"ats"`
}
