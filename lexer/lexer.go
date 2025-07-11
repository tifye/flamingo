package lexer

import (
	"strings"
	"unicode/utf8"

	"github.com/tifye/flamingo/assert"
	"github.com/tifye/flamingo/token"
)

const (
	eof rune = -1
)

type stateFunc func(*Lexer) stateFunc

type Lexer struct {
	input  string
	tokens chan token.Token
	state  stateFunc
	start  int
	pos    int
	width  int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		// Channel must be large enough to support the largest
		// amount of tokens that can be outputed from a single state
		tokens: make(chan token.Token, 6),
		state:  lexText,
	}
	return l
}

func (l *Lexer) NextToken() token.Token {
	for {
		select {
		case item := <-l.tokens:
			return item
		default:
			l.state = l.state(l)
		}
	}
}

func (l *Lexer) emit(typ token.TokenType) {
	if typ == token.EOF {
		tok := token.Token{Type: typ}
		l.tokens <- tok
		l.start = l.pos
		return
	}

	assert.Assert(l.pos > l.start, "pos must be past start")

	tok := token.Token{
		Type:    typ,
		Literal: l.input[l.start:l.pos],
	}
	l.tokens <- tok
	l.start = l.pos
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = size
	l.pos += l.width
	return r
}

func (l *Lexer) backup() {
	l.pos -= l.width
	assert.Assert(l.pos >= 0, "pos must be larger than or equal to zero")
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *Lexer) runUntil(valid string) {
	for !strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func lexText(l *Lexer) stateFunc {
	assert.AssertNotNil(l)

	for {
		if strings.HasPrefix(l.input[l.pos:], "<") {
			// There might not be any text
			if l.pos > l.start {
				l.emit(token.TEXT)
			}
			return lexTag
		}

		if l.next() == eof {
			if l.pos > l.start {
				l.emit(token.TEXT)
			}
			break
		}
	}

	l.emit(token.EOF)
	return nil
}

func lexTag(l *Lexer) stateFunc {
	assert.AssertNotNil(l)

	ch := l.next()
	assert.Assert(ch == '<', "expected '<'")
	l.emit(token.LEFT_CHEV)

	if l.accept("/") {
		l.emit(token.SLASH)
	}

	l.runUntil("> ")
	ch = l.peek()
	if ch == eof {
		if l.pos > l.start {
			l.emit(token.IDENT)
		}

		l.emit(token.EOF)
		return nil
	}

	l.emit(token.IDENT)

	// todo: check all types of empty characters
	l.acceptRun(" ")

	if l.accept(">") {
		l.emit(token.RIGHT_CHEV)
		return lexText
	}

	return lexAttribute
}

func lexAttribute(l *Lexer) stateFunc {
	assert.Assert(!l.accept(" "), "expected no empty characters")

	l.runUntil("=")
	if l.peek() == eof {
		if l.pos > l.start {
			l.emit(token.IDENT)
		}
		l.emit(token.EOF)
		return nil
	}

	l.emit(token.IDENT)

	if !l.accept("=") {
		l.emit(token.ERROR)
		return nil
	}
	l.emit(token.ASSIGN)

	if !l.accept(`"`) {
		l.emit(token.ERROR)
		return nil
	}
	l.emit(token.QUOTE)

	l.runUntil(`"`)
	if l.peek() == eof {
		if l.pos > l.start {
			l.emit(token.TEXT)
		}
		l.emit(token.EOF)
		return nil
	}

	l.emit(token.TEXT)

	if !l.accept(`"`) {
		l.emit(token.ERROR)
		return nil
	}
	l.emit(token.QUOTE)

	l.acceptRun(" ")

	if l.accept(">") {
		l.emit(token.RIGHT_CHEV)
		return lexText
	}

	return lexAttribute
}
