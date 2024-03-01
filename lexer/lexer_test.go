package lexer

import (
	"interingo/token"
	"testing"
)

type expectedReturn struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func test(t *testing.T, input string, tests []expectedReturn, testsName string) {
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("testsName[%s] - tests[%d] - tokentype wrong. expected=%q, got=%q",
				testsName, i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("testsName[%s] - tests[%d] - literal wrong. expected=%q, got=%q",
				testsName, i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func testOperator1(t *testing.T) {
	input := `=+(){},;`
	tests := []expectedReturn{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	test(t, input, tests, "Operators 1")
}

func testKeywords1(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
    x + y;
};`
	tests := []expectedReturn{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	test(t, input, tests, "Keywords 1")
}

func testKeywords2(t *testing.T) {
	input := `if (5 < 10) {
    return true;
} else {
    return false;
}`
	tests := []expectedReturn{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}
	test(t, input, tests, "Keywords 2")
}

func testOperator2(t *testing.T) {
	input := `!-/*5;
5 < 10 > 5;`
	tests := []expectedReturn{
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	test(t, input, tests, "Operator 2")
}

func testOperator3(t *testing.T) {
	input := `10 == 10;
10 != 9;`
	tests := []expectedReturn{
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	test(t, input, tests, "Operator 3")
}

func TestNextToken(t *testing.T) {
	testOperator1(t)
	testOperator2(t)
	testOperator3(t)
	testKeywords1(t)
	testKeywords2(t)
	testBinding1(t)
}

func testBinding1(t *testing.T) {
	input := "let x = 5;"
	input += "\nlet y = 10;"
	input += "\nlet foobar = add(5, 5);"
	input += "let barfoo = 5 * 5 / 10 + 18 - add(5, 5) + multiply(124);"
	tests := []expectedReturn{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "foobar"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.COMMA, ","},
		{token.INT, "5"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.LET, "let" },
		{token.IDENT, "barfoo"},
		{token.ASSIGN , "="},
		{token.INT , "5"},
		{token.ASTERISK , "*"},
		{token.INT , "5"},
		{token.SLASH , "/"},
		{token.INT , "10"},
		{token.PLUS , "+"},
		{token.INT , "18"},
		{token.MINUS , "-"},
		{token.IDENT , "add"},
		{token.LPAREN , "("},
		{token.INT , "5"},
		{token.COMMA, ","},
		{token.INT , "5"},
		{token.RPAREN , ")"},
		{token.PLUS , "+"},
		{token.IDENT , "multiply"},
		{token.LPAREN , "("},
		{token.INT , "124"},
		{token.RPAREN , ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	test(t, input, tests, "Binding 1")
}
