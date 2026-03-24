package parser

import (
	"interingo/pkg/share"
)

type ParserError struct {
	Message string
	Range   share.Range
}
