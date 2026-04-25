// This handle core runtime that evaluate the code
// - Move over core language handling into this file
// - REPL - Eval will using this instead

package runtime

import (
	"fmt"
	"maps"
	"sort"
	"strings"

	"interingo/pkg/evaluator"
	"interingo/pkg/lexer"
	"interingo/pkg/object"
	"interingo/pkg/parser"
	"interingo/pkg/runtime/embed"
	"interingo/pkg/runtime/native"
	"interingo/pkg/token"
)

type RuntimeType string

const (
	EMBED  = RuntimeType("embed")
	NATIVE = RuntimeType("native")
)

type RuntimeState int

const (
	RUNTIME_UNKNOWN RuntimeState = iota // 0
	RUNTIME_START
	RUNTIME_END
)

type Core struct {
	Env              *object.Environment
	Verbose          bool
	lifeCycleHandler LifeCycleHandler // Only work with embed type
	state            RuntimeState
}

func (c *Core) loadNativeBuiltIn() {
	c.Env.Set("exit", &native.SystemExit{
		Environment: c.Env,
	})
}

func (c *Core) loadEmbedBuiltIn() {
	c.Env.Set("exit", &embed.SystemExit{
		Environment: c.Env,
	})
}

func (c *Core) State() RuntimeState {
	return c.state
}

func NewCore(t RuntimeType, lifeCycleHandler LifeCycleHandler) *Core {
	env := object.NewEnvironment()

	c := &Core{
		Env:              env,
		Verbose:          false,
		lifeCycleHandler: lifeCycleHandler,
		state:            RUNTIME_START,
	}

	if t == NATIVE {
		c.loadNativeBuiltIn()
	} else {
		c.loadEmbedBuiltIn()
	}

	// Allow toggle verbose to update verbose flag of this runtime core
	c.Env.Set("toggleVerbose", &ToggleVerbose{
		Core: c,
	})
	c.Env.Set("getRuntimeInfo", &GetRuntimeInfo{
		env:  env,
		Core: c,
	})
	c.Env.Set("help", &Help{
		env:  env,
		Core: c,
	})

	return c
}

func (c *Core) Eval(req EvalRequest) (*EvalResponseSuccess, *EvalResponseError, *VerboseInfo) {
	if c.state == RUNTIME_END {
		return nil, &EvalResponseError{
			Error: fmt.Errorf("Runtime state already ended!"),
		}, nil
	}

	l := lexer.New(req.Data)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors) != 0 {
		error := &EvalResponseError{ParserErrors: p.Errors}
		// Verbose should work even parser error
		var verboseInfo *VerboseInfo
		if c.Verbose {
			verboseInfo = getVerboseInfomation(l, p)
		}

		return nil, error, verboseInfo
	}

	evaluated := evaluator.Eval(program, c.Env)
	// This here as Eval can change runtime verbose state
	var verboseInfo *VerboseInfo
	if c.Verbose {
		verboseInfo = getVerboseInfomation(l, p)
	}

	// Check if it is a system exit
	if sysExit, ok := evaluated.(*object.SystemExit); ok {
		err := &EvalResponseError{
			SystemExit: &SystemExit{
				Code: sysExit.Code,
			},
		}
		c.state = RUNTIME_END
		if c.lifeCycleHandler != nil {
			c.lifeCycleHandler.Exit()
		}
		return nil, err, verboseInfo
	}

	result := &EvalResponseSuccess{
		Output: nil,
	}

	if evaluated != nil {
		output := evaluated.Inspect()
		result.Output = &output
	}

	return result, nil, verboseInfo
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

// Return help infomation, which contain a single string to provide to the User
func (c *Core) Help() string {
	infos := c.Env.GetAllBuiltinInfos()
	// Make sure the keys reponse in specific order
	keys := []string{}
	for k, _ := range infos {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var helpInfo strings.Builder

	for i, k := range keys {
		if i > 0 {
			fmt.Fprintf(&helpInfo, "\n")
		}
		fmt.Fprintf(&helpInfo, "\t- %s(%s) : %s", k, object.FnParamsInspect(infos[k].Parameters()), infos[k].Description())
	}
	return helpInfo.String()
}

// Return help infomation, which contain a single string to provide to the User
func (c *Core) GetRuntimeInfo() string {
	infos := c.Env.GetAllStoreData()
	return infos
}
