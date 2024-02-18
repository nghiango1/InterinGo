package repl

import (
	"io"
	"log"
	"main/evaluator"
	"main/lexer"
	"main/parser"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

const EVAL_UNDEFINE = "Seem eval function not implemented yet"

func Start(in io.Reader, out io.Writer) {
	signalCapture(in, out)
}

func handle(line string, out io.Writer) {
	l := lexer.New(line)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		errorsHandler(p.Errors(), out)
		return
	}
	evaluated := evaluator.Eval(program)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	} else {
		errs := []string{EVAL_UNDEFINE, "Your statements are: ", program.String()}
		errorsHandler(errs, out)
	}
}

func errorsHandler(errors []string, out io.Writer) {
	io.WriteString(out, "Parse got "+strconv.Itoa(len(errors))+" errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func usage(w io.Writer) {
	io.WriteString(w, "Language keywords:\n")
	io.WriteString(w, "\tlet\n")
	io.WriteString(w, "\tif\n")
	io.WriteString(w, "Built-in commands:\n")
	io.WriteString(w, "exit()\n")
	io.WriteString(w, "help()\n")
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("let"),
	readline.PcItem("if ("),
	readline.PcItem("exit()"),
	readline.PcItem("help()"),
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
		Prompt:          ">>> ",
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
		switch {
		case line == "help()":
			usage(l.Stderr())
		case line == "exit()":
			goto exit
		case line == "":
		default:
			handle(line, out)
		}
	}
exit:
}
