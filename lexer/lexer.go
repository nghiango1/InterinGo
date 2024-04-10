package lexer

import (
	"interingo/share"
	"interingo/token"
)

type Lexer struct {
	input        string
	position     int                     // current position in input (points to current char)
	readPosition int                     // current reading position in input (after current char)
	ch           byte                    // current char under examination
	SkipedChar   int                     // skiped white-space - for Verbose mode
	SkipedLine   int                     // skiped comment line - for Verbose mode
	TokenCount   map[token.TokenType]int // count all token - for Verbose mode
	Line         int                     // current position in input (line - 0 index)
	Character    int                     // current position in input (position in line - 0 index)
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	if share.VerboseMode {
		l.TokenCount = make(map[token.TokenType]int)
	}
	l.readChar()
	return l
}

func (l *Lexer) peakChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.Character += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		if share.VerboseMode {
			l.SkipedChar += 1
		}
		l.readChar()
	}
}

func (l *Lexer) skipCurrentLine() string {
	if share.VerboseMode {
		l.SkipedLine += 1
	}

	pos := l.position
	// Read until end of line or stop when reach EOF
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
		// This goes pass the current token
		l.readChar()

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
		l.readChar()
	case '(':
		tok = newToken(token.LPAREN, l.ch)
		l.readChar()
	case '-':
		tok = newToken(token.MINUS, l.ch)
		l.readChar()

	case '!':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
		l.readChar()

	case '*':
		tok = newToken(token.ASTERISK, l.ch)
		l.readChar()
	case '/':
		if l.peakChar() == '/' {
			tok.Type = token.COMMENT
			literal := l.skipCurrentLine()
			tok.Literal = literal
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
		l.readChar()

	case '>':
		tok = newToken(token.GT, l.ch)
		l.readChar()
	case '<':
		tok = newToken(token.LT, l.ch)
		l.readChar()
	case ')':
		tok = newToken(token.RPAREN, l.ch)
		l.readChar()
	case ',':
		tok = newToken(token.COMMA, l.ch)
		l.readChar()
	case '+':
		tok = newToken(token.PLUS, l.ch)
		l.readChar()
	case '{':
		tok = newToken(token.LBRACE, l.ch)
		l.readChar()
	case '}':
		tok = newToken(token.RBRACE, l.ch)
		l.readChar()

	case '\n':
		tok = newToken(token.EOL, l.ch)
		l.readChar()
		l.Line += 1
		l.Character = 0

	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
		l.readChar()

	case '\r':
		if l.peakChar() == '\n' {
			ch := l.ch
			l.readChar()
			literal := ch + l.ch
			tok = newToken(token.EOL, literal)
		} else {
			tok = newToken(token.EOL, l.ch)
		}
		l.readChar()
		l.Line += 1
		l.Character = 0

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		} else if isDigit(l.ch) {
			tok.Literal = l.readDigit()
			tok.Type = token.INT
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	if share.VerboseMode {
		l.TokenCount[tok.Type] += 1
	}
	return tok
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readDigit() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readString() string {
	l.readChar()
	pos := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[pos:l.position]
}


func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
