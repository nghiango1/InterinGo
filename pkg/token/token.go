// token/token.go
package token

import (
	"fmt"
	"interingo/pkg/share"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Start   share.Position
	End     share.Position
}

// token/token.go
const (
	COMMENT = TokenType("COMMENT")
	EOL     = TokenType("EOL")
	ILLEGAL = TokenType("ILLEGAL")
	EOF     = TokenType("EOF")
	// Identifiers + literals
	IDENT = TokenType("IDENT") // add, foobar, x, y, ...
	INT   = TokenType("INT")   // 1343456
	// Operators
	ASSIGN   = TokenType("=")
	PLUS     = TokenType("+")
	MINUS    = TokenType("-")
	BANG     = TokenType("!")
	ASTERISK = TokenType("*")
	SLASH    = TokenType("/")
	GT       = TokenType(">")
	LT       = TokenType("<")
	EQ       = TokenType("==")
	NOT_EQ   = TokenType("!=")
	// Delimiters
	COMMA     = TokenType(",")
	SEMICOLON = TokenType(";")
	LPAREN    = TokenType("(")
	RPAREN    = TokenType(")")
	LBRACE    = TokenType("{")
	RBRACE    = TokenType("}")
	// Keywords
	FUNCTION = TokenType("FN")
	LET      = TokenType("LET")
	IF       = TokenType("IF")
	ELSE     = TokenType("ELSE")
	RETURN   = TokenType("RETURN")
	TRUE     = TokenType("TRUE")
	FALSE    = TokenType("FALSE")
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t *Token) String() string {
	return fmt.Sprintf("Token { %v, %v, %v }", t.Type, t.Start, t.Literal)
}
