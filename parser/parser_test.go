package parser

import (
	source "go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tifye/flamingo/lexer"
)

func TestNestedSimpleComponents(t *testing.T) {
	input := `<div class="bg-rose-500">mino</div>`
	fset := source.NewFileSet()
	f := fset.AddFile("", fset.Base(), len(input))
	l := lexer.NewLexer(f, input)
	p := NewParser(l)

	root := p.Parse()
	noParserErrors(t, p)

	require.NotNil(t, root)
	require.NotNil(t, root.Fragment, "expected a Fragment node")
	assert.Equal(t, 1, len(root.Fragment.Nodes), "expected component to contain 1 node")
}

func noParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if assert.Empty(t, errs, "expected no errors") {
		return
	}

	for _, msg := range errs {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
