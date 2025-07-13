package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
)

func TestNestedSimpleComponents(t *testing.T) {
	input := `<div class="bg-rose-500">mino</div>`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	root := p.Parse()
	require.NotNil(t, root)
	require.NotNil(t, root.Fragment, "expected a Fragment node")
	assert.Equal(t, 1, len(root.Fragment.Nodes), "expected component to contain 1 node")

	require.Equal(t, "div", root.Fragment.Nodes[0].TokenLit())
	div := root.Fragment.Nodes[0].(*ast.Component)
	assert.Equal(t, "mino", div.Nodes[0].TokenLit())
}
