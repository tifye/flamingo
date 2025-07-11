package token

import "fmt"

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
