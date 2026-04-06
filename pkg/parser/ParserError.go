package parser

import (
	"interingo/pkg/share"
)

type ParserError struct {
	Message string      `json:"message"`
	Range   share.Range `json:"range"`
}
