package token

import (
	"fmt"

	"github.com/tifye/flamingo/assert"
)

type TokenType int

const (
	ERROR TokenType = iota
	EOF

	LEFT_CHEV
	RIGHT_CHEV
	SLASH

	IDENT
	ASSIGN
	QUOTE
	COLON
	ON

	TEXT
	GO_EXPR
	GO_CODE
)

var ttStr = map[TokenType]string{
	ERROR:      "ERROR",
	EOF:        "EOF",
	LEFT_CHEV:  "LEFT_CHEV",
	RIGHT_CHEV: "RIGHT_CHEV",
	SLASH:      "SLASH",
	IDENT:      "IDENT",
	ASSIGN:     "ASSIGN",
	QUOTE:      "QUOTE",
	COLON:      "COLON",
	ON:         "ON",
	TEXT:       "TEXT",
	GO_EXPR:    "GO_EXPR",
	GO_CODE:    "GO_CODE",
}

func (tt TokenType) String() string {
	str, ok := ttStr[tt]
	assert.Assert(ok, fmt.Sprintf("missing TokenType entry for %d", tt))
	return str
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	switch t.Type {
	case EOF:
		return "EOF"
	case ERROR:
		return "ERROR"
	}
	if len(t.Literal) > 10 {
		return fmt.Sprintf("%.10q...", t.Literal)
	}
	return fmt.Sprintf("%q", t.Literal)
}
