package parser

import (
	"fmt"
	"interingo/pkg/ast"
	"interingo/pkg/lexer"
	"interingo/pkg/share"
	"interingo/pkg/token"
	"log/slog"
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
	Program        *ast.Program
	Lexer          *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	DocumentTokens []DocumentToken
	Errors         []ParserError
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
			p.handlerNextToken()
		} else {
			break
		}
	}
}

// Assumming we handle parse infix, the last insert value kind can have updated
// Semantic kind after we understand more on the context
// We have to go back for some value
func (p *Parser) reverseIndentityLiteralKind(literal string, newKind SemanticTokenType) error {
	for i := len(p.DocumentTokens) - 1; i >= 0; i-- {
		if p.DocumentTokens[i].Token.Type != token.IDENT {
			continue
		}
		if p.DocumentTokens[i].Token.Literal != literal {
			continue
		}
		p.DocumentTokens[i].Kind = newKind
		return nil
	}

	return fmt.Errorf("Can't find token match `%v` profile", literal)
}

// This just bind semantic kind without consider of the context, we rely on reverse
// back into DocumentTokens previous value to cover correct semantic kind
func (p *Parser) handlerNextToken() {
	if p.curToken.Type != "" {
		p.DocumentTokens = append(p.DocumentTokens, DocumentTokenWrap(p.curToken, TokenTypeToSemanticKind(p.curToken.Type)))
	}
	p.curToken = p.peekToken
	p.peekToken = p.Lexer.NextToken()
}

func (p *Parser) nextToken() {
	p.handlerNextToken()
	p.skipExtras()
}

func (p *Parser) ParseProgram() *ast.Program {
	if p.Program != nil {
		return p.Program
	}

	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	slog.Debug(fmt.Sprintf("At: %v", p.curToken))
	program.Range.End = p.curToken.End

	for _, curToken := range p.DocumentTokens {
		if curToken.Kind == SemanticTokenTypeComment {
			program.Comments = append(program.Comments, curToken.Unwrap())
		}
	}
	p.Program = program
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
	stmt.Range.Start = p.curToken.Start

	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression != nil {
		stmt.Range.End = stmt.Expression.GetRange().End
	}

	if p.peekTokenIs(token.SEMICOLON) {
		stmt.Range.Start = p.curToken.End
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("Expression expect, but %s found", t)
	p.Errors = append(p.Errors, ParserError{
		Message: msg,
		Range: share.Range{
			Start: p.curToken.Start,
			End:   p.curToken.End,
		},
	})
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
	exp.Range.Start = p.curToken.Start

	exp.Function = function
	if exp.Function != nil {
		p.reverseIndentityLiteralKind(exp.Function.TokenLiteral(), SemanticTokenTypeFunction)
	}

	for !p.peekTokenIs(token.RPAREN) {
		p.nextToken() // Skip the '(' and ',' token
		exp.Arguments = append(exp.Arguments, p.parseExpression(LOWEST))
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	exp.Range.End = p.curToken.End
	p.nextToken() // Skip the ')' token

	if p.peekTokenIs(token.SEMICOLON) {
		exp.Range.End = p.curToken.End
		p.nextToken()
	}
	return exp
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	stmt.Range.Start = p.curToken.Start
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)
	if stmt.ReturnValue != nil {
		stmt.Range.End = stmt.ReturnValue.GetRange().End
	} else {
		stmt.Range.End = stmt.Token.End
	}

	if p.peekTokenIs(token.SEMICOLON) {
		stmt.Range.End = p.curToken.End
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
	block.Range.Start = p.curToken.Start

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		curr := p.parseStatement()
		if curr != nil {
			block.Statements = append(block.Statements, curr)
		}
		p.nextToken()
		block.Range.End = p.curToken.End
	}

	return block
}

func isFunctionLiteral(expr ast.Expression) bool {
	_, ok := expr.(*ast.FunctionLiteral)
	return ok
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}
	stmt.Range.Start = p.curToken.Start

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
		Range: share.Range{
			Start: p.curToken.Start,
			End:   p.curToken.End,
		},
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	value := p.parseExpression(LOWEST)
	stmt.Value = value

	if isFunctionLiteral(value) {
		p.reverseIndentityLiteralKind(stmt.Name.Value, SemanticTokenTypeFunction)
	} else {
		p.reverseIndentityLiteralKind(stmt.Name.Value, SemanticTokenTypeVariable)
	}

	if value != nil {
		stmt.Range.End = value.GetRange().End
	}

	if p.peekTokenIs(token.SEMICOLON) {
		stmt.Range.End = p.curToken.End
		p.nextToken()
	}

	return stmt
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{Lexer: l}
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
	function.Range.Start = p.curToken.Start

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
	if function.Body != nil {
		function.Range.End = function.Body.GetRange().End
	}
	if p.peekTokenIs(token.SEMICOLON) {
		function.Range.End = p.curToken.End
		p.nextToken()
	}

	return function
}

func (p *Parser) parseIfElseExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}
	expression.Range.Start = p.curToken.Start

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
	if expression.Consequence != nil {
		expression.Range.End = expression.Consequence.GetRange().End
	}
	expression.Range.End = expression.Consequence.GetRange().End

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
		if expression.Alternative != nil {
			expression.Range.End = expression.Alternative.GetRange().End
		}
	}

	if p.peekTokenIs(token.SEMICOLON) {
		expression.Range.End = p.curToken.End
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
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
		Range: share.Range{
			Start: p.curToken.Start,
			End:   p.curToken.End,
		},
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
		Range: share.Range{
			Start: p.curToken.Start,
			End:   p.curToken.End,
		},
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	lit.Range.Start = p.curToken.Start

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer",
			p.curToken.Literal)
		p.Errors = append(p.Errors, ParserError{
			Message: msg,
			Range: share.Range{
				Start: p.curToken.Start,
				End:   p.curToken.End,
			}})
		return nil
	}
	lit.Value = value

	lit.Range.End = p.curToken.End
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	expression.Range.Start = p.curToken.Start
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	expression.Range.End = p.curToken.End
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	expression.Range.Start = left.GetRange().Start
	precedences := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedences)

	if expression.Right != nil {
		expression.Range.End = expression.Right.GetRange().End
	}
	return expression
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.Errors = append(p.Errors, ParserError{
		Message: msg,
		Range: share.Range{
			Start: p.peekToken.Start,
			End:   p.peekToken.End,
		},
	})
}
