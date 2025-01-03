package lexer

import "gomonkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = newTokenFromChar(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newTokenFromChar(token.PLUS, l.ch)
	case '-':
		tok = newTokenFromChar(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = newTokenFromChar(token.BANG, l.ch)
		}
	case '/':
		tok = newTokenFromChar(token.SLASH, l.ch)
	case '*':
		tok = newTokenFromChar(token.ASTERISK, l.ch)
	case '<':
		tok = newTokenFromChar(token.LT, l.ch)
	case '>':
		tok = newTokenFromChar(token.GT, l.ch)
	case ';':
		tok = newTokenFromChar(token.SEMICOLON, l.ch)
	case ',':
		tok = newTokenFromChar(token.COMMA, l.ch)
	case '(':
		tok = newTokenFromChar(token.LPAREN, l.ch)
	case ')':
		tok = newTokenFromChar(token.RPAREN, l.ch)
	case '{':
		tok = newTokenFromChar(token.LBRACE, l.ch)
	case '}':
		tok = newTokenFromChar(token.RBRACE, l.ch)
	case '[':
		tok = newTokenFromChar(token.LBRACKET, l.ch)
	case ']':
		tok = newTokenFromChar(token.RBRACKET, l.ch)
	case '"':
		tok = newToken(token.STRING, l.readString())
	case 0:
		tok = newToken(token.EOF, "")
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok = newToken(token.LookupIdent(literal), literal)
			return tok
		} else if isDigit(l.ch) {
			tok = newToken(token.INT, l.readNumber())
			return tok
		} else {
			tok = newTokenFromChar(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func newTokenFromChar(tokenType token.TokenType, ch byte) token.Token {
	return newToken(tokenType, string(ch))
}
