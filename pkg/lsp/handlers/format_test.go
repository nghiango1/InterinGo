package handlers

import (
	"encoding/json"
	"fmt"
	"interingo/pkg/lexer"
	"interingo/pkg/parser"
	"interingo/pkg/share"
	"os"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestCommentLineFormat(t *testing.T) {
	input := `//This is a coment line
add(1, 2 * 3, 4 + 5);
// This is a second commented line`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	jsonData, err := json.Marshal(program)
	if err != nil {
		fmt.Printf("[ERROR] Got: %v\n", err)
	} else {
		fmt.Printf("[INFO] Data: (%v) %v\n", len(jsonData), string(jsonData))
	}

	fmt.Println(FormatedAST(program, protocol.FormattingOptions{}, 0))
}

func TestFunctionFormat(t *testing.T) {
	input := `	// 5
let identity = fn(x) {
		// 5
    return x;
	// 5
};
identity(5);`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	jsonData, err := json.Marshal(program)
	if err != nil {
		fmt.Printf("[ERROR] Got: %v\n", err)
	} else {
		fmt.Printf("[INFO] Data: (%v) %v\n", len(jsonData), string(jsonData))
	}

	fmt.Println(FormatedAST(program, protocol.FormattingOptions{}, 0))
}

func TestMain(m *testing.M) {
	share.SetDefaultLog()
    os.Exit(m.Run())
}
