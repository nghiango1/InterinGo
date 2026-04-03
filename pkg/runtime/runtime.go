// This handle core runtime that evaluate the code
// - Move over core language handling into this file
// - REPL - Eval will using this instead

package runtime

import (
	"interingo/pkg/evaluator"
	"interingo/pkg/lexer"
	"interingo/pkg/object"
	"interingo/pkg/parser"
	"interingo/pkg/runtime/native"
	"interingo/pkg/share"
	"interingo/pkg/token"
	"maps"
)

type RuntimeType string

const (
	EMBED  = RuntimeType("embed")
	NATIVE = RuntimeType("native")
)

type Core struct {
	Env     object.Environment
	Verbose bool
}

func (c *Core) loadNativeBuiltIn(env *object.Environment) {
	env.Set("exit", &native.SystemExit{
		Environment: env,
	})
}

func (c *Core) loadEmbedBuiltIn(env *object.Environment) {
	// Should not be able to exit
	// env.Set("exit", &embed.SystemExit{ Environment: env, })
}

func NewCore(t RuntimeType) *Core {
	env := object.NewEnvironment()

	c := &Core{
		Env:     *env,
		Verbose: share.VerboseMode,
	}

	if t == NATIVE {
		c.loadNativeBuiltIn(env)
	} else {
		c.loadEmbedBuiltIn(env)
	}

	return c
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

	evaluated := evaluator.Eval(program, &c.Env)
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
