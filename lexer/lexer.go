package lexer

import (
	"fmt"
	"pythia/token"
	"strings"
)

type Lexer struct {
	input        []rune
	position     int  // 현재 문자의 위치
	readPosition int  // 현재 문자의 다음
	ch           rune // 현재 읽고 있는 문자
}

func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) GetErrorInfo() string {
	row := 0
	col := 0
	chars := len(l.input)
	i := 0

	for i < l.readPosition && i < chars {
		if l.input[i] == '\n' {
			row++
			col = 0
		}

		i++
		col++
	}
	return fmt.Sprintf("got %c at line %d, around %d", l.ch, row, col)
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		tok = l.makeTwoCharToken(l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = l.makeTwoCharToken(l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '%':
		tok = newToken(token.PERCENT, l.ch)
	case '<':
		tok = l.makeTwoCharToken(l.ch)
	case '>':
		tok = l.makeTwoCharToken(l.ch)
	case '&':
		tok = l.makeTwoCharToken(l.ch)
	case '|':
		tok = l.makeTwoCharToken(l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '.':
		tok.Type = token.DOT
		tok.Literal = l.readInstruction()
	case rune(0):
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			return l.newNumberToken()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = rune(0)
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) newNumberToken() token.Token {
	var tok token.Token

	tok.Literal = l.readNumber()
	if strings.Contains(tok.Literal, ".") {
		tok.Type = token.FLOAT
	} else {
		tok.Type = token.INT
	}

	return tok
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return rune(0)
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}

	}

	return string(l.input[position:l.position])
}

func (l *Lexer) readInstruction() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == 0 {
			break
		}
	}

	return string(l.input[position:l.position])
}

func (l *Lexer) makeTwoCharToken(currChar rune) token.Token {
	var tok token.Token

	switch currChar {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			literal := string(currChar) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, currChar)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			literal := string(currChar) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, currChar)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			literal := string(currChar) + string(l.ch)
			tok = token.Token{Type: token.LT_OR_EQ, Literal: literal}
		} else {
			tok = newToken(token.LT, currChar)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			literal := string(currChar) + string(l.ch)
			tok = token.Token{Type: token.GT_OR_EQ, Literal: literal}
		} else {
			tok = newToken(token.GT, currChar)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			literal := string(currChar) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal}
		} else {
			tok = newToken(token.ILLEGAL, currChar)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			literal := string(currChar) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal}
		} else {
			tok = newToken(token.ILLEGAL, currChar)
		}
	}

	return tok
}
