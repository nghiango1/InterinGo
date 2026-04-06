// This handle read-input, evaluation, print-output, loop
// Using readline to handle stdin/out only

package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"interingo/pkg/parser"
	"interingo/pkg/runtime"

	"github.com/chzyer/readline"
)

const EVAL_UNDEFINE = "Seem eval function not implemented yet"

type Repl struct {
	core *runtime.Core
	in   io.Reader
	out  io.Writer
}

func NewRepl(evalCore *runtime.Core, in io.Reader, out io.Writer) *Repl {
	if evalCore == nil {
		evalCore = runtime.NewCore(runtime.NATIVE, nil)
	}

	res := &Repl{
		core: evalCore,
		in:   in,
		out:  out,
	}

	evalCore.Env.Set(
		"print", &Print{
			env: res.core.Env,
			r:   res,
		},
	)

	return res
}

func (r *Repl) Start() {
	r.signalCapture()
}

func (r *Repl) Handle(line string) error {
	if line == "" {
		return nil
	}

	result, error, verbose := r.core.Eval(runtime.EvalRequest{
		Data: line,
	})

	if verbose != nil {
		r.printVerboseInfomation(verbose)
	}

	if error != nil {
		if len(error.ParserErrors) != 0 {
			r.parsingErrorsHandler(error.ParserErrors)
			// Already handle
			// return fmt.Errorf("Parser error")
			return nil
		}

		if error.Error != nil {
			return error.Error
		}

		if error.SystemExit != nil {
			return fmt.Errorf("System exit called")
		}

		return fmt.Errorf("Unknown error")
	}

	if result == nil {
		return fmt.Errorf("Evaluation result not found error")
	}

	if result.Output != nil {
		io.WriteString(r.out, *result.Output)
		io.WriteString(r.out, "\n")
	}

	return nil
}

func (r *Repl) printVerboseInfomation(info *runtime.VerboseInfo) {
	data, err := json.MarshalIndent(info, "> ", "    ")
	if err != nil {
		return
	}
	r.out.Write(data)
	io.WriteString(r.out, "\n")
}

func (r *Repl) parsingErrorsHandler(errors []parser.ParserError) {
	io.WriteString(r.out, "Errors when parsing:\n")
	for _, e := range errors {
		io.WriteString(r.out, "\t"+e.Message+"\n")
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (r *Repl) signalCapture() {
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("let "),
		readline.PcItem("if ("),
	)

	for k, v := range r.core.Env.GetAllBuiltinInfos() {
		funcCall := "("
		if len(v.Parameters().Standard) == 0 && len(v.Parameters().Default) == 0 && !v.Parameters().Rest {
			funcCall += ")"
		}
		completer.Children = append(completer.Children, readline.PcItem(k + funcCall))
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:          ">> ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		err = r.Handle(line)
		if err != nil {
			panic(err)
		}
	}
}
