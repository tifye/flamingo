package parser

import (
	source "go/token"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tifye/flamingo/ast"
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

func TestParseElement(t *testing.T) {
	t.Run("<div></div>", func(t *testing.T) {
		input := `<div></div>`
		el, err := ParseElement(input)
		assert.NoError(t, err)
		assert.NotNil(t, el)
		ast.Inspect(el, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.Element:
				assert.Equal(t, "div", n.Name.Name)
				assert.Equal(t, 1, int(n.LeftChevron))
				assert.Equal(t, 11, int(n.RightChevron))
			}
			return true
		})
	})

	t.Run("<div><meep><div></div></meep><mino></mino></div>", func(t *testing.T) {
		input := `<div><meep><div></div></meep><mino></mino></div>`
		elements := []string{"div", "div", "meep", "mino"}
		el, err := ParseElement(input)
		assert.NoError(t, err)
		assert.NotNil(t, el)
		ast.Inspect(el, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.Element:
				idx := slices.Index(elements, n.Name.Name)
				elements = slices.Delete(elements, idx, idx+1)
			}
			return true
		})
		assert.Empty(t, elements)
	})

	t.Run("<self />", func(t *testing.T) {
		input := `<self />`
		el, err := ParseElement(input)
		assert.NoError(t, err)
		assert.NotNil(t, el)
		ast.Inspect(el, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.Element:
				assert.Equal(t, "self", n.Name.Name)
				assert.Equal(t, 1, int(n.LeftChevron))
				assert.Equal(t, 8, int(n.RightChevron))
			}
			return true
		})
	})
}

func TestCodeBlock(t *testing.T) {
	input := `---
func _() {}
---
<test></test>
`
	fset := source.NewFileSet()
	el, err := ParseFile(fset, "", input)
	assert.NoError(t, err)
	assert.NotNil(t, el)
	didParseCodeBlock := false
	ast.Inspect(el, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CodeBlock:
			assert.Equal(t, "func _() {}\n", n.Code)
			assert.Equal(t, source.Pos(1), n.Pos())
			assert.Equal(t, source.Pos(17), n.End())
			didParseCodeBlock = true
		}
		return true
	})
	assert.True(t, didParseCodeBlock, "expected to parse code block")
}

func TestAttribute(t *testing.T) {
	t.Run(`empty string literal`, func(t *testing.T) {
		input := `<test isTrue=""/>`
		el, err := ParseElement(input)
		assert.NoError(t, err)
		assert.NotNil(t, el)
		ast.Inspect(el, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.Attribute:
				assert.Equal(t, "isTrue", n.Name.Name)
				assert.Equal(t, "", n.ValueLiteral)
			}
			return true
		})
	})

	t.Run(`just attribute name (boolean)`, func(t *testing.T) {
		input := `<test isTrue/>`
		el, err := ParseElement(input)
		assert.NoError(t, err)
		assert.NotNil(t, el)
		ast.Inspect(el, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.Attribute:
				assert.Equal(t, "isTrue", n.Name.Name)
				assert.Equal(t, "true", n.ValueLiteral)
			}
			return true
		})
	})

	t.Run(`multiline attributes`, func(t *testing.T) {
		input := `<test
			meep="meep"
		/>`
		el, err := ParseElement(input)
		assert.NoError(t, err)
		assert.NotNil(t, el)
		ast.Inspect(el, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.Attribute:
				assert.Equal(t, "meep", n.Name.Name)
				assert.Equal(t, "meep", n.ValueLiteral)
			}
			return true
		})
	})
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
