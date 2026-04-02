// This handle how to evaluate the code
// - Move over core language handling into this file
// - REPL - Eval will using this instead

package evaluator

import (
	"interingo/pkg/lexer"
	"interingo/pkg/object"
	"interingo/pkg/parser"
	"interingo/pkg/share"
	"interingo/pkg/token"
	"maps"
)

type Core struct {
	Env     object.Environment
	Verbose bool
}

func loadBuiltIn(env *object.Environment) {
	env.Set("exit", &object.SystemExit{
		Environment: env,
	})
}

func NewCore() *Core {
	env := object.NewEnvironment()
	loadBuiltIn(env)

	return &Core{
		Env:     *env,
		Verbose: share.VerboseMode,
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

	if len(p.Errors) != 0 {
		error := &EvalResponseError{ParserErrors: p.Errors}
		return nil, error, verbose
	}

	evaluated := Eval(program, &c.Env)
	var err *EvalResponseError

	result := &EvalResponseSuccess{
		Output: nil,
	}

	if evaluated != nil {
		output := evaluated.Inspect()
		result.Output = &output
	}

	return result, err, verbose
}

func (c *Core) ToggleVerbose() {
	c.Verbose = !c.Verbose
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
