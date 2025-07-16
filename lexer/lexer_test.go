package lexer

import (
	gtoken "go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tifye/flamingo/token"
)

func TestNextToken(t *testing.T) {
	fset := gtoken.NewFileSet()
	t.Run(`<div class="p-4">mino</div>`, func(t *testing.T) {
		input := `<div class="p-4">mino</div>`
		f := fset.AddFile("", fset.Base(), len(input))

		tests := []struct {
			expectedType token.TokenType
		}{
			{token.LEFT_CHEV},
			{token.IDENT},
			{token.IDENT},
			{token.ASSIGN},
			{token.QUOTE},
			{token.TEXT},
			{token.QUOTE},
			{token.RIGHT_CHEV},
			{token.TEXT},
			{token.LEFT_CHEV},
			{token.SLASH},
			{token.IDENT},
			{token.RIGHT_CHEV},
			{token.EOF},
		}

		l := NewLexer(f, input)
		for i, tt := range tests {
			tok := l.NextToken()
			assert.Equal(t, tt.expectedType, tok.Type, "Token idx %d", i)
		}
	})

	t.Run(`<div>mino<div>`, func(t *testing.T) {
		input := `<div>mino</div>`
		f := fset.AddFile("", fset.Base(), len(input))

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

		l := NewLexer(f, input)
		for i, tt := range tests {
			tok := l.NextToken()
			assert.Equal(t, tt.expectedType, tok.Type, "Token idx %d", i)
		}
	})

	t.Run(`<div><span>Meep</span></div>`, func(t *testing.T) {
		input := `<div><span>Meep</span></div>`
		f := fset.AddFile("", fset.Base(), len(input))

		tests := []struct {
			expectedType token.TokenType
		}{
			{token.LEFT_CHEV},
			{token.IDENT},
			{token.RIGHT_CHEV},
			{token.LEFT_CHEV},
			{token.IDENT},
			{token.RIGHT_CHEV},
			{token.TEXT},
			{token.LEFT_CHEV},
			{token.SLASH},
			{token.IDENT},
			{token.RIGHT_CHEV},
			{token.LEFT_CHEV},
			{token.SLASH},
			{token.IDENT},
			{token.RIGHT_CHEV},
			{token.EOF},
		}

		l := NewLexer(f, input)
		for i, tt := range tests {
			tok := l.NextToken()
			assert.Equal(t, tt.expectedType, tok.Type, "Token idx %d", i)
		}
	})
}

func TestTxtToken(t *testing.T) {
	input := `<div>`
	fset := gtoken.NewFileSet()
	f := fset.AddFile("", fset.Base(), len(input))

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.LEFT_CHEV},
		{token.IDENT},
		{token.RIGHT_CHEV},
		{token.EOF},
	}

	l := NewLexer(f, input)
	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "Token idx %d", i)
	}
}
