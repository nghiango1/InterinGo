package parser

import (
	"main/ast"
	"main/lexer"
	"main/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two token no curToken and peekToken is not null
	p.nextToken()
	p.nextToken()
	return p
}
