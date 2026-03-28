// This handle how to evaluate the code
// - Move over core language handling into this file
// - REPL - Eval will using this instead

package core

import (
	"interingo/pkg/evaluator"
	"interingo/pkg/lexer"
	"interingo/pkg/object"
	"interingo/pkg/parser"
	"interingo/pkg/token"
	"maps"
)

type Core struct {
	Env     object.Environment
	Verbose bool
}

func NewCore() *Core {
	return &Core{
		Env: *object.NewEnvironment(),
	}
}

func (c *Core) Eval(req EvalRequest) (*EvalResponseSuccess, *EvalResponseError, *VerboseInfo) {
	l := lexer.New(req.Data)
	p := parser.New(l)
	program := p.ParseProgram()

	var verbose *VerboseInfo
	if c.Verbose {
		verbose = getVerboseInfomation(l, p)
	}

	if len(p.Errors()) != 0 {
		error := &EvalResponseError{ParserErrors: p.Errors()}
		return nil, error, verbose
	}

	evaluated := evaluator.Eval(program, &c.Env)

	result := &EvalResponseSuccess{
		Output: nil,
	}

	if evaluated != nil {
		output := evaluated.Inspect()
		result.Output = &output
	}

	return result, nil, verbose
}

func getVerboseInfomation(l *lexer.Lexer, p *parser.Parser) *VerboseInfo {
	var result VerboseInfo
	result.Lexer.WhitespaceSkip = l.SkipedChar
	result.Lexer.CommentLine = l.SkipedCommentLine
	result.Lexer.Token = make(map[token.TokenType]int)

	maps.Copy(result.Lexer.Token, l.TokenCount)

	result.Parser.Ats = p.Program
	return &result
}
