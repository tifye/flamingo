package staticweb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicRende(t *testing.T) {
	builder := strings.Builder{}
	renderer := NewStaticWebRenderer(&builder)
	expected := `<div title="test"><span>Izu</span></div>`

	parent := renderer.NewComponent("div")
	parent.SetAttribute("title", "test")

	child := renderer.NewComponent("span")
	child.SetAttribute("innerText", "Izu")
	renderer.Append(parent, child)

	renderer.Render(parent)

	result := builder.String()
	assert.Equal(t, expected, result)
}
