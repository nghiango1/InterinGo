package parser

import (
	"main/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let y = 10;
let foobar = 8383838;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Parse Program function return nul")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program stagements is less than 3, got len=%d ", len(program.Statements))
	}
	tests = []struct{
		expectedIdentifier string
	} {
		{"x"}
		{"y"}
		{"foobar"}
	}

	for i, tt := range tests {
		stmt
	}
}

