package core

import (
	"interingo/pkg/evaluator"
	"interingo/pkg/lexer"
	"interingo/pkg/object"
	"interingo/pkg/parser"
	"interingo/pkg/share"
)

type Core struct {
	Env object.Environment
}

func NewCore() *Core {
	return &Core{
		Env: *object.NewEnvironment(),
	}
}

func (c *Core) Eval(req EvaluateRequest) (EvaluateResponse, error) {
	l := lexer.New(req.Data)
	p := parser.New(l)
	program := p.ParseProgram()

	if share.VerboseMode {
		r.printVerboseInfomation(l, p)
	}

	if len(p.Errors()) != 0 {
		errorsHandler(p.Errors(), r.out)
		return
	}

	evaluated := evaluator.Eval(program, &c.Env)
	return EvaluateResponse{
		Output: evaluated.Inspect(),
	}, nil
}
