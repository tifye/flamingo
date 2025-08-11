package lexer

import (
	source "go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tifye/flamingo/token"
)

func TestNextToken(t *testing.T) {
	fset := source.NewFileSet()
	t.Run(`<div class="p-4">mino</div>`, func(t *testing.T) {
		input := `<div class="p-4">mino</div>`
		f := fset.AddFile("", fset.Base(), len(input))

		tests := []struct {
			expectedType token.TokenType
		}{
			{token.LEFT_CHEVRON},
			{token.IDENT},
			{token.IDENT},
			{token.ASSIGN},
			{token.QUOTE},
			{token.TEXT},
			{token.QUOTE},
			{token.RIGHT_CHEVRON},
			{token.TEXT},
			{token.LEFT_CHEVRON},
			{token.SLASH},
			{token.IDENT},
			{token.RIGHT_CHEVRON},
			{token.EOF},
		}

		l := NewLexer(f, input)
		for i, tt := range tests {
			tok := l.NextToken()
			assert.Equal(t, tt.expectedType.String(), tok.Type.String(), "Token idx %d", i)
		}
	})

	t.Run(`<div>mino<div>`, func(t *testing.T) {
		input := `<div>mino</div>`
		f := fset.AddFile("", fset.Base(), len(input))

		tests := []struct {
			expectedType token.TokenType
		}{
			{token.LEFT_CHEVRON},
			{token.IDENT},
			{token.RIGHT_CHEVRON},
			{token.TEXT},
			{token.LEFT_CHEVRON},
			{token.SLASH},
			{token.IDENT},
			{token.RIGHT_CHEVRON},
			{token.EOF},
		}

		l := NewLexer(f, input)
		for i, tt := range tests {
			tok := l.NextToken()
			assert.Equal(t, tt.expectedType.String(), tok.Type.String(), "Token idx %d", i)
		}
	})

	t.Run(`<div><span>Meep</span></div>`, func(t *testing.T) {
		input := `<div><span>Meep</span></div>`
		f := fset.AddFile("", fset.Base(), len(input))

		tests := []struct {
			expectedType token.TokenType
		}{
			{token.LEFT_CHEVRON},
			{token.IDENT},
			{token.RIGHT_CHEVRON},
			{token.LEFT_CHEVRON},
			{token.IDENT},
			{token.RIGHT_CHEVRON},
			{token.TEXT},
			{token.LEFT_CHEVRON},
			{token.SLASH},
			{token.IDENT},
			{token.RIGHT_CHEVRON},
			{token.LEFT_CHEVRON},
			{token.SLASH},
			{token.IDENT},
			{token.RIGHT_CHEVRON},
			{token.EOF},
		}

		l := NewLexer(f, input)
		for i, tt := range tests {
			tok := l.NextToken()
			assert.Equal(t, tt.expectedType.String(), tok.Type.String(), "Token idx %d", i)
		}
	})
}

func TestTxtToken(t *testing.T) {
	input := `
<div class="meep">
	mino
	meep
</div>`
	fset := source.NewFileSet()
	f := fset.AddFile("", fset.Base(), len(input))

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.LEFT_CHEVRON},
		{token.IDENT},
		{token.IDENT},
		{token.ASSIGN},
		{token.QUOTE},
		{token.TEXT},
		{token.QUOTE},
		{token.RIGHT_CHEVRON},
		{token.TEXT},
		{token.LEFT_CHEVRON},
		{token.SLASH},
		{token.IDENT},
		{token.RIGHT_CHEVRON},
		{token.EOF},
	}

	l := NewLexer(f, input)
	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "Token idx %d", i)
	}
}

func TestCodeBlock(t *testing.T) {
	input := `
---
func _() {}
---
<test></test>
`
	fset := source.NewFileSet()
	f := fset.AddFile("", fset.Base(), len(input))
	tests := []struct {
		expectedType token.TokenType
	}{
		{token.CODE_FENCE},
		{token.GO_CODE},
		{token.CODE_FENCE},
		{token.LEFT_CHEVRON},
		{token.IDENT},
		{token.RIGHT_CHEVRON},
		{token.LEFT_CHEVRON},
		{token.SLASH},
		{token.IDENT},
		{token.RIGHT_CHEVRON},
		{token.EOF},
	}

	l := NewLexer(f, input)
	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "Token idx %d, expected %s, got %s", i, tt.expectedType, tok.Type)
	}
}
