package token

import (
	"fmt"
	source "go/token"
)

//go:generate stringer -type=TokenType
type TokenType int

type Pos = source.Pos

const (
	ERROR TokenType = iota
	EOF

	LEFT_CHEVRON
	RIGHT_CHEVRON
	SLASH

	IDENT
	ASSIGN
	QUOTE
	COLON
	ON

	TEXT
	GO_EXPRESSION
	GO_CODE
	CODE_FENCE
)

type Token struct {
	Pos     Pos
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
