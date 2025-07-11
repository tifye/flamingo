package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tifye/flamingo/token"
)

func TestNextToken(t *testing.T) {
	input := `<div>mino</div>`

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.LEFT_CHEV},
		{token.IDENT},
		{token.RIGHT_CHEV},
		{token.TEXT},
		{token.LEFT_CHEV},
		{token.SLASH},
		{token.IDENT},
		{token.RIGHT_CHEV},
		{token.EOF},
	}

	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "Token idx %d", i)
	}
}
