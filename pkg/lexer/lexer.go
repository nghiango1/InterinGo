package lexer

import (
	"interingo/pkg/share"
	"interingo/pkg/token"
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

func (l *Lexer) Position() token.Position {
	return token.Position{
		Line:      l.Line,
		Character: l.Character - 1,
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	start := l.Position()
	switch l.ch {
	case '=':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = newToken(token.EQ, literal, start)
		} else {
			tok = newToken(token.ASSIGN, string(l.ch), start)
		}
	case ';':
		tok = newToken(token.SEMICOLON, string(l.ch), start)
	case '(':
		tok = newToken(token.LPAREN, string(l.ch), start)
	case '-':
		tok = newToken(token.MINUS, string(l.ch), start)
	case '!':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = newToken(token.NOT_EQ, literal, start)
		} else {
			tok = newToken(token.BANG, string(l.ch), start)
		}
	case '*':
		tok = newToken(token.ASTERISK, string(l.ch), start)
	case '/':
		if l.peakChar() == '/' {
			tok.Type = token.COMMENT
			literal := l.skipCurrentLine()
			tok.Literal = literal

			return tok
		} else {
			tok = newToken(token.SLASH, string(l.ch), start)
		}
	case '>':
		tok = newToken(token.GT, string(l.ch), start)
	case '<':
		tok = newToken(token.LT, string(l.ch), start)
	case ')':
		tok = newToken(token.RPAREN, string(l.ch), start)
	case ',':
		tok = newToken(token.COMMA, string(l.ch), start)
	case '+':
		tok = newToken(token.PLUS, string(l.ch), start)
	case '{':
		tok = newToken(token.LBRACE, string(l.ch), start)
	case '}':
		tok = newToken(token.RBRACE, string(l.ch), start)
	case '\n':
		tok = newToken(token.EOL, string(l.ch), start)
		l.Line += 1
		l.Character = 0
	case '\r':
		if l.peakChar() == '\n' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = newToken(token.EOL, literal, start)
		} else {
			tok = newToken(token.EOL, string(l.ch), start)
		}
		l.Line += 1
		l.Character = 0
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Start = start
			tok.End = l.Position()

			if share.VerboseMode {
				l.TokenCount[tok.Type] += 1
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readDigit()
			tok.Type = token.INT
			tok.Start = start
			tok.End = l.Position()

			if share.VerboseMode {
				l.TokenCount[tok.Type] += 1
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, string(l.ch), start)
		}
	}
	l.readChar()

	if share.VerboseMode {
		l.TokenCount[tok.Type] += 1
	}

	tok.End = token.Position{
		Line:      l.Line,
		Character: l.Character,
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

func newToken(tokenType token.TokenType, ch string, start token.Position) token.Token {
	return token.Token{Type: tokenType, Literal: ch, Start: start}
}
