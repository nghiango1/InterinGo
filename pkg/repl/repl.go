package repl

import (
	"fmt"
	"io"
	"log"
	"strings"

	"interingo/pkg/lexer"
	"interingo/pkg/object"
	"interingo/pkg/parser"
	"interingo/pkg/service/core"
	"interingo/pkg/share"

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
		share.VerboseMode = !share.VerboseMode
		if share.VerboseMode {
			io.WriteString(r.out, "Verbose mode enable\n")
		} else {
			io.WriteString(r.out, "Verbose mode disable\n")
		}
	case line == "":
	default:
		r.codeHandle(line)
	}
}

func (r *Repl) printVerboseInfomation(l *lexer.Lexer, p *parser.Parser) {
	lexerVerbose := fmt.Sprintf("Lexer information:\n\tSkip whitespace = %d\n\tSkip comment line = %d\n\tToken found:\n", l.SkipedChar, l.SkipedLine)
	io.WriteString(r.out, lexerVerbose)

	for k, v := range l.TokenCount {
		io.WriteString(r.out, fmt.Sprintf("\t\t%v: %d\n", k, v))
	}

	// TO DO: Print parser infomation, Print full AST tree presentation
	parseVerbose := "Parse infomation:\n\tProgram statement parsed:\n"
	io.WriteString(r.out, parseVerbose)
	for _, v := range p.Program.Statements {
		io.WriteString(r.out, fmt.Sprintf("\t\t%T: %v\n", v, v.String()))
	}
}

func (r *Repl) codeHandle(line string) {
	if line == "" {
		return
	}

	evaluated, err := r.core.Eval(core.EvaluateRequest{
		Data: line,
	})

	if err != nil {
		r.errorsHandler([]string{err.Error()})
	}

	if share.VerboseMode {
		r.printVerboseInfomation(l, p)
	}

	if len(p.Errors()) != 0 {
		errorsHandler(p.Errors(), r.out)
		return
	}

	if evaluated != nil {
		io.WriteString(r.out, evaluated.Inspect())
		io.WriteString(r.out, "\n")
	}
}

func (r *Repl) errorsHandler(errors []string) {
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
