package parser

import (
	"fmt"
	"interingo/ast"
	"interingo/lexer"
	"interingo/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

// This map the operator token to their piority class
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// This check what is the next token piority
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// This check what is the current token piority
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) skipExtras() {
	for p.curToken.Type != token.EOF {
		if p.curToken.Type == token.COMMENT || p.curToken.Type == token.EOL {
			p.curToken = p.peekToken
			p.peekToken = p.l.NextToken()
		} else {
			break
		}
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	p.skipExtras()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token: p.curToken,
	}

	exp.Function = function

	for !p.peekTokenIs(token.RPAREN) {
		p.nextToken() // Skip the '(' and ',' token
		exp.Arguments = append(exp.Arguments, p.parseExpression(LOWEST))
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	p.nextToken() // Skip the ')' token

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return exp
}
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if !p.curTokenIs(token.LBRACE) {
		return nil
	}

	block := &ast.BlockStatement{
		Token: p.curToken,
	}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		curr := p.parseStatement()
		if curr != nil {
			block.Statements = append(block.Statements, curr)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	value := p.parseExpression(LOWEST)
	stmt.Value = value

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfElseExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// Read two token no curToken and peekToken is not null
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	function := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		ident := p.parseIdentifier().(*ast.Identifier)
		function.Parameters = append(function.Parameters, ident)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	function.Body = p.parseBlockStatement()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return function
}

func (p *Parser) parseIfElseExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer",
			p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedences := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedences)
	return expression
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
