package ast

import (
	"bytes"
	"interingo/pkg/token"
)

type Range struct {
	Start token.Position `json:"start"`
	End   token.Position `json:"end"`
}

type Node interface {
	TokenLiteral() string
	String() string
	GetRange() Range
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement   `json:"statements,omitempty"`
	Comments   []token.Token `json:"comments,omitempty"`
	Range      Range         `json:"range"`
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for i, s := range p.Statements {
		if i > 0 {
			out.WriteString("; ")
		}
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) GetRange() Range {
	return Range{
		Start: p.Range.Start,
		End:   p.Range.End,
	}
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier `json:"name"`
	Value Expression  `json:"value"`
	Range Range       `json:"range"`
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) GetRange() Range {
	return Range{
		Start: ls.Range.Start,
		End:   ls.Range.End,
	}
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
	Range       Range
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) GetRange() Range {
	return Range{
		Start: rs.Range.Start,
		End:   rs.Range.End,
	}
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(rs.ReturnValue.String())

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the Expression
	Expression Expression
	Range      Range
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) GetRange() Range {
	return Range{
		Start: es.Range.Start,
		End:   es.Range.End,
	}
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
	Range Range
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) GetRange() Range {
	return Range{
		Start: i.Range.Start,
		End:   i.Range.End,
	}
}

func (i *Identifier) String() string {
	return i.Value
}

type InfixExpression struct {
	Token    token.Token // The in token, e.g. + - * /
	Operator string
	Left     Expression
	Right    Expression
	Range    Range
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) GetRange() Range {
	return Range{
		Start: ie.Range.Start,
		End:   ie.Range.End,
	}
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
	Range    Range
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) GetRange() Range {
	return Range{
		Start: pe.Range.Start,
		End:   pe.Range.End,
	}
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
	Range Range
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) GetRange() Range {
	return Range{
		Start: il.Range.Start,
		End:   il.Range.End,
	}
}
func (il *IntegerLiteral) String() string { return il.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
	Range Range
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) GetRange() Range {
	return Range{
		Start: b.Range.Start,
		End:   b.Range.End,
	}
}
func (b *Boolean) String() string { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token //the "if" token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Range       Range
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) GetRange() Range {
	return Range{
		Start: ie.Range.Start,
		End:   ie.Range.End,
	}
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token //the "{" token
	Parameters []*Identifier
	Body       *BlockStatement
	Range      Range
}

func (fe *FunctionLiteral) expressionNode()      {}
func (fe *FunctionLiteral) TokenLiteral() string { return fe.Token.Literal }
func (fe *FunctionLiteral) GetRange() Range {
	return Range{
		Start: fe.Range.Start,
		End:   fe.Range.End,
	}
}
func (fe *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("fn")
	out.WriteString("(")
	for i, ident := range fe.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(ident.String())
	}
	out.WriteString(") ")
	out.WriteString(fe.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token //the "(" token
	Function  Expression
	Arguments []Expression
	Range     Range
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) GetRange() Range {
	return Range{
		Start: ce.Range.Start,
		End:   ce.Range.End,
	}
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, exp := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(exp.String())
	}
	out.WriteString(")")
	return out.String()
}

type BlockStatement struct {
	Token      token.Token //the "{" token
	Statements []Statement
	Range      Range
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) GetRange() Range {
	return Range{
		Start: bs.Range.Start,
		End:   bs.Range.End,
	}
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")

	for i, stmt := range bs.Statements {
		if i > 0 {
			out.WriteString("; ")
		}
		out.WriteString(stmt.String())
	}
	out.WriteString(" }")
	return out.String()
}
