package parser

import (
	"interingo/pkg/token"
)

type DocumentToken struct {
	Token token.Token
	Kind  SemanticTokenType
}

func TokenTypeToSemanticKind(t token.TokenType) SemanticTokenType {
	var tokenType SemanticTokenType
	switch t {
	case token.COMMENT:
		tokenType = SemanticTokenTypeComment
	case
		token.ASSIGN,
		token.PLUS,
		token.MINUS,
		token.BANG,
		token.ASTERISK,
		token.SLASH,
		token.GT,
		token.LT,
		token.EQ,
		token.NOT_EQ,
		token.COMMA,
		token.SEMICOLON,
		token.LPAREN,
		token.RPAREN,
		token.LBRACE,
		token.RBRACE:
		tokenType = SemanticTokenTypeOperator
	case
		token.IF,
		token.LET,
		token.ELSE,
		token.RETURN,
	    token.FUNCTION:
		tokenType = SemanticTokenTypeKeyword
	case
		token.TRUE,
		token.FALSE:
		tokenType = SemanticTokenTypeMacro
	case token.INT:
		tokenType = SemanticTokenTypeNumber
	case token.IDENT:
		tokenType = SemanticTokenTypeVariable
	case token.ILLEGAL:
		tokenType = SemanticTokenTypeType
	default:
		tokenType = SemanticTokenTypeType
	}
	return tokenType
}

func TokenToSemanticKind(t token.Token) SemanticTokenType {
	return TokenTypeToSemanticKind(t.Type)
}

func DocumentTokenWrap(token token.Token, kind SemanticTokenType) DocumentToken {
	return DocumentToken{
		Token: token,
		Kind:  kind,
	}
}

func DocumentTokenUnwrap(documentToken DocumentToken) token.Token {
	return documentToken.Token
}

func (documentToken *DocumentToken) Unwrap() token.Token {
	return documentToken.Token
}
