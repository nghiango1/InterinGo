package repl

import (
	"fmt"
	"io"
	"log"
	"interingo/ast"
	"interingo/evaluator"
	"interingo/lexer"
	"interingo/object"
	"interingo/parser"
	"interingo/share"
	"strings"

	"github.com/chzyer/readline"
)

const EVAL_UNDEFINE = "Seem eval function not implemented yet"

var env *object.Environment = nil

func Start(in io.Reader, out io.Writer) {
	signalCapture(in, out)
}

func Handle(line string, out io.Writer) {
	if env == nil {
		env = object.NewEnvironment()
	}

	switch {
	case line == "help()":
		usage(out)
	case line == "exit()":
		io.WriteString(out, "exit() only work in REPL CLI session")
	case line == "toggleVerbose()":
		share.VerboseMode = !share.VerboseMode
		if share.VerboseMode {
			io.WriteString(out, "Verbose mode enable")
		} else {
			io.WriteString(out, "Verbose mode disable")
		}
	case line == "":
	default:
		codeHandle(line, out, env)
	}
}

func printVerboseInfomation(l *lexer.Lexer, p *parser.Parser, program *ast.Program, out io.Writer) {
	lexerVerbose := fmt.Sprintf("Lexer information:\n\tSkip whitespace = %d\n\tSkip comment line = %d\n\tToken found:\n", l.SkipedChar, l.SkipedLine)
	io.WriteString(out, lexerVerbose)

	for k, v := range l.TokenCount {
		io.WriteString(out, fmt.Sprintf("\t\t%v: %d\n", k, v))
	}

	// TO DO: Print parser infomation, Print full AST tree presentation
	parseVerbose := fmt.Sprintf("Parse infomation:\n\tProgram statement parsed:\n")
	io.WriteString(out, parseVerbose)
	for _, v := range program.Statements {
		io.WriteString(out, fmt.Sprintf("\t\t%T: %v\n", v, v.String()))
	}
}

func codeHandle(line string, out io.Writer, env *object.Environment) {
	if line == "" {
		return
	}
	l := lexer.New(line)
	p := parser.New(l)
	program := p.ParseProgram()

	if share.VerboseMode {
		printVerboseInfomation(l, p, program, out)
	}

	if len(p.Errors()) != 0 {
		errorsHandler(p.Errors(), out)
		return
	}

	if share.VerboseMode {
		io.WriteString(out, "Evaluation result:\n\n")
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func errorsHandler(errors []string, out io.Writer) {
	io.WriteString(out, "Errors when parsing:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
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

func signalCapture(in io.Reader, out io.Writer) {
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
			Handle(line, out)
		}
	}
}
