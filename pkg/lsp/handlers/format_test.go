package handlers

import (
	"encoding/json"
	"fmt"
	"interingo/pkg/lexer"
	"interingo/pkg/parser"
	"interingo/pkg/share"
	"log/slog"
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

func TestFunctionFormatHard(t *testing.T) {
	input := `// 5
let c = fn(x)  {
// hw
    return x;
    //asd'
}(12,2);
// return //as das/ahehehe  a/ a/ a a/a /a a/ a aa/ a 
        // as das ssdas dahuhuhuhuh u hu hu uh uh  h  uh h h h taht tso oos oso sso os os 
// return hehe(hhh, h, hh) // asasdasdasd
        // let c =
// if (a(c){a{c}/3 * 4 - a */} /)
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	jsonData, err := json.Marshal(program)
	if err != nil {
		slog.Error(fmt.Sprintf("Got: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Data: (%v) %v\n", len(jsonData), string(jsonData)))
	}

	fmt.Println(FormatedAST(program, protocol.FormattingOptions{}, 0))
}

func TestMain(m *testing.M) {
	share.SetDefaultLog()
	os.Exit(m.Run())
}
