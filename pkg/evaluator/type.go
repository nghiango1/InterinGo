package evaluator

import (
	"encoding/json"
	"interingo/pkg/parser"
	"interingo/pkg/token"
)

type EvalRequest struct {
	Data string `json:"data"`
}

// Might have no Output return, eg: `let x = 2`
type EvalResponseSuccess struct {
	Output *string `json:"output,omitempty"`
}

func (resp *EvalResponseSuccess) String() string {
	b, err := json.Marshal(resp)
	if err != nil {
		return ""
	}
	return string(b)
}

type EvalResponseError struct {
	ParserErrors []parser.ParserError `json:"parserErrors,omitempty"`
}

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

func (vi *VerboseInfo) String() (string, error) {
	data, err := json.MarshalIndent(vi, "> ", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
