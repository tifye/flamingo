package staticweb

import (
	"fmt"
	"io"

	"github.com/tifye/flamingo/render"
)

type Component struct {
	name  string
	text  string
	comps []render.Component
	attrs map[string]any
}

func (c *Component) SetAttribute(key string, val any) {
	if c.attrs == nil {
		c.attrs = make(map[string]any)
	}
	if key == "innerText" {
		c.text = fmt.Sprintf("%s", val)
		return
	}
	c.attrs[key] = val
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) Children() []render.Component {
	return c.comps
}

func (c *Component) Attributes() map[string]any {
	return c.attrs
}

type StaticWebRenderer struct {
	writer io.StringWriter
}

func NewStaticWebRenderer(writer io.StringWriter) *StaticWebRenderer {
	return &StaticWebRenderer{
		writer: writer,
	}
}

func (r *StaticWebRenderer) NewComponent(name string) render.Component {
	return &Component{
		name: name,
	}
}

func (r StaticWebRenderer) Render(comps ...render.Component) {
	for _, c := range comps {
		comp, ok := c.(*Component)
		if !ok {
			panic("invalid comp type")
		}

		r.writeComponent(comp)
	}
}

func (r *StaticWebRenderer) writeComponent(c *Component) {
	_, _ = r.writer.WriteString("<")
	_, _ = r.writer.WriteString(c.name)

	for key, val := range c.attrs {
		_, _ = r.writer.WriteString(fmt.Sprintf(` %s="%s"`, key, val))
	}

	_, _ = r.writer.WriteString(">")

	if c.text != "" {
		_, _ = r.writer.WriteString(c.text)
	}

	for _, cc := range c.comps {
		comp, ok := cc.(*Component)
		if !ok {
			panic("invalid comp type")
		}

		r.writeComponent(comp)
	}

	_, _ = r.writer.WriteString("</")
	_, _ = r.writer.WriteString(c.name)
	_, _ = r.writer.WriteString(">")
}

func (r *StaticWebRenderer) Append(parent render.Component, child render.Component) {
	parentc, ok := parent.(*Component)
	if !ok {
		panic("invalid comp type")
	}

	childc, ok := child.(*Component)
	if !ok {
		panic("invalid comp type")
	}

	parentc.comps = append(parentc.comps, childc)
}
