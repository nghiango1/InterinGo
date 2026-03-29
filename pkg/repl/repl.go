// This handle read-input, evaluation, print-output, loop
// Using readline to handle stdin/out only

package repl

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"interingo/pkg/evaluator"
	"interingo/pkg/object"
	"interingo/pkg/parser"

	"github.com/chzyer/readline"
)

const EVAL_UNDEFINE = "Seem eval function not implemented yet"

type Repl struct {
	core evaluator.Core
	in   io.Reader
	out  io.Writer
}

func NewRepl(evalCore *evaluator.Core, in io.Reader, out io.Writer) *Repl {
	if evalCore == nil {
		evalCore = evaluator.NewCore()
	}

	return &Repl{
		core: *evalCore,
		in:   in,
		out:  out,
	}
}

func (r *Repl) Start() {
	r.signalCapture()
}

func (r *Repl) Handle(input string) {
	switch input {
	case "help()":
		usage(r.out)
	case "exit()":
		io.WriteString(r.out, "exit() only work in REPL CLI session, but let me reset all variable for you\n")
		r.core.Env = *object.NewEnvironment()
	case "toggleVerbose()":
		r.core.ToggleVerbose()
		if r.core.Verbose {
			io.WriteString(r.out, "Verbose mode enable\n")
		} else {
			io.WriteString(r.out, "Verbose mode disable\n")
		}
	case "":
	default:
		r.codeHandle(input)
	}
}

func (r *Repl) printVerboseInfomation(info *evaluator.VerboseInfo) {
	data, err := json.MarshalIndent(info, "> ", "    ")
	if err != nil {
		return
	}
	r.out.Write(data)
	io.WriteString(r.out, "\n")
}

func (r *Repl) codeHandle(line string) {
	if line == "" {
		return
	}

	result, error, verbose := r.core.Eval(evaluator.EvalRequest{
		Data: line,
	})

	if verbose != nil {
		r.printVerboseInfomation(verbose)
	}

	if error != nil {
		if len(error.ParserErrors) != 0 {
			r.parsingErrorsHandler(error.ParserErrors)
		}
	}

	if result != nil {
		if result.Output != nil {
			io.WriteString(r.out, *result.Output)
			io.WriteString(r.out, "\n")
		}
	}
}

func (r *Repl) parsingErrorsHandler(errors []parser.ParserError) {
	io.WriteString(r.out, "Errors when parsing:\n")
	for _, e := range errors {
		io.WriteString(r.out, "\t"+e.Message+"\n")
	}
}

func usage(w io.Writer) {
	io.WriteString(w, "Built-in commands:\n")
	io.WriteString(w, "- toggleVerbose(): Toggle verbose mode - print more infomation about Lexer, Parse and Evaluator\n")
	io.WriteString(w, "- help(): Print this help\n")
	io.WriteString(w, "- exit(): End this REPL session\n")
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("let"),
	readline.PcItem("if ("),
	readline.PcItem("exit()"),
	readline.PcItem("help()"),
	readline.PcItem("toggleVerbose()"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (r *Repl) signalCapture() {
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
		if line == "exit()" {
			break
		} else {
			r.Handle(line)
		}
	}
}
