package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"interingo/pkg/object"
	"interingo/pkg/service/core"

	"github.com/chzyer/readline"
)

const EVAL_UNDEFINE = "Seem eval function not implemented yet"

type Repl struct {
	core core.Core
	in   io.Reader
	out  io.Writer
}

func NewRepl(in io.Reader, out io.Writer) *Repl {
	return &Repl{
		core: *core.NewCore(),
		in:   in,
		out:  out,
	}
}

func (r *Repl) Start() {
	r.signalCapture()
}

func (r *Repl) Handle(line string) {
	switch {
	case line == "help()":
		usage(r.out)
	case line == "exit()":
		io.WriteString(r.out, "exit() only work in REPL CLI session, but let me reset all variable for you\n")
		r.core.Env = *object.NewEnvironment()
	case line == "toggleVerbose()":
		r.core.Verbose = !r.core.Verbose
		if r.core.Verbose {
			io.WriteString(r.out, "Verbose mode enable\n")
		} else {
			io.WriteString(r.out, "Verbose mode disable\n")
		}
	case line == "":
	default:
		r.codeHandle(line)
	}
}

func (r *Repl) printVerboseInfomation(info *core.VerboseInfo) {
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

	evaluateResult, error := r.core.Eval(core.EvaluateRequest{
		Data: line,
	})

	if error != nil {
		if e, ok := error.(*core.ParserErrorResponse); ok {
			if r.core.Verbose {
				r.printVerboseInfomation(e.Verbose)
			}
			r.parsingErrorsHandler(e.Errors)
		} else {
			io.WriteString(r.out, fmt.Sprintf("Unknown errors when handling code, got: %v\n", error))
		}
	}

	if evaluateResult != nil {
		if r.core.Verbose {
			r.printVerboseInfomation(evaluateResult.Verbose)
		}

		io.WriteString(r.out, evaluateResult.Output)
		io.WriteString(r.out, "\n")
	}
}

func (r *Repl) parsingErrorsHandler(errors []string) {
	io.WriteString(r.out, "Errors when parsing:\n")
	for _, msg := range errors {
		io.WriteString(r.out, "\t"+msg+"\n")
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
