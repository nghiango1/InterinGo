package repl

import (
	"bufio"
	"fmt"
	"io"
	"main/lexer"
	"main/parser"
	// "main/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// var tok token.Token
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		parser := parser.New(l)
		program := parser.ParseProgram()
		stmts := program.Statements
		errors := parser.Errors()
		if len(errors) != 0 {
			fmt.Errorf("Parser has %d errors", len(errors))
			for _, msg := range errors {
				fmt.Errorf("Parser error: %q", msg)
			}
		} else {
			for i := 0; i < len(stmts); i++ {
				fmt.Fprintf(out, "Statement %s\n", stmts[i].String())
			}
		}
		// for tok = l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// 	//fmt.Fprintf(out, "%+v\n", tok)
		// 	fmt.Fprintf(out, "{token.%v, \"%s\"},\n", tok.Type, tok.Literal)
		// }
		// fmt.Fprintf(out, "{token.%v, \"%s\"},\n", tok.Type, tok.Literal)
	}
}
