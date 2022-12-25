package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

var EOF rune = rune(0)

type Lexer struct {
	reader *bufio.Reader
}

type TokenType int

// Valid tokens
const (
	TokenTypeUnknown TokenType = iota
	TokenTypeIdent
	TokenTypeConfig
	TokenTypeEndStmt
	TokenTypeSpace
	TokenTypeHyphen
	TokenTypeEOF
)

type Token struct {
	Type  TokenType
	Value string
}

func NewLexer(r io.Reader) *Lexer {
	lexer := &Lexer{
		reader: bufio.NewReader(r),
	}
	return lexer
}

func (l *Lexer) read() rune {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		return EOF
	}
	return ch
}

func (l *Lexer) unread() {
	l.reader.UnreadRune()
}

func (l *Lexer) Scan() Token {
	ch := l.read()

	switch ch {
	case EOF:
		return Token{Type: TokenTypeEOF}
	case '+':
		return Token{Type: TokenTypeConfig, Value: string(ch)}
	case '\n':
		return Token{Type: TokenTypeEndStmt, Value: string(ch)}
	case ' ':
		l.unread()
		return l.scanRune(' ', TokenTypeSpace)
	case '-':
		l.unread()
		return l.scanRune('-', TokenTypeHyphen)
	default:
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			l.unread()
			return l.scanIdent()
		}
		return Token{Type: TokenTypeUnknown}
	}
}

func (l *Lexer) scanRune(acceptedRune rune, tokenType TokenType) Token {
	ch := l.read()
	value := &strings.Builder{}
	value.WriteRune(ch)
	for {
		ch := l.read()
		if ch == EOF {
			break
		} else if ch != acceptedRune {
			l.unread()
			break
		} else {
			value.WriteRune(ch)
		}
	}
	return Token{Type: tokenType, Value: value.String()}
}

func (l *Lexer) scanIdent() Token {
	ch := l.read()
	value := &strings.Builder{}
	value.WriteRune(ch)

	for {
		ch := l.read()
		if ch == EOF {
			break
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			l.unread()
			break
		} else {
			value.WriteRune(ch)
		}
	}
	return Token{Type: TokenTypeIdent, Value: value.String()}
}
