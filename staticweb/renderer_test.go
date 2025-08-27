package staticweb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicRende(t *testing.T) {
	builder := strings.Builder{}
	renderer := NewStaticWebRenderer(&builder)
	expected := `<div title="test"><span></span></div>`

	comp := renderer.NewComponent("div")
	comp.SetAttribute("title", "test")
	renderer.Append(comp, renderer.NewComponent("span"))
	renderer.Render(comp)

	result := builder.String()
	assert.Equal(t, expected, result)
}
