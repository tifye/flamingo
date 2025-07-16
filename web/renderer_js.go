package web

import (
	"log"
	"strings"
	"syscall/js"

	"github.com/tifye/flamingo/render"
)

type jsFunc func(this js.Value, args []js.Value) any

type WebComponent struct {
	name  string
	text  string
	comps []render.Component
	attrs map[string]any
	el    *js.Value
}

func (c *WebComponent) Element() *js.Value {
	return c.el
}

func (c *WebComponent) SetAttribute(key string, val any) {
	if c.attrs == nil {
		c.attrs = make(map[string]any)
	}
	c.attrs[key] = val

	if c.el != nil {
		switch key {
		case "innerText", "value":
			c.el.Set(key, val)
		default:
			c.el.Call("setAttribute", key, val)
		}
	}
}

func (c *WebComponent) Name() string {
	return c.name
}

func (c *WebComponent) Children() []render.Component {
	return c.comps
}

func (c *WebComponent) Attributes() map[string]any {
	return c.attrs
}

type DOMRenderer struct {
	doc  js.Value
	body js.Value
}

func NewDOMRenderer() *DOMRenderer {
	doc := js.Global().Get("document")
	body := doc.Get("body")

	return &DOMRenderer{
		doc:  doc,
		body: body,
	}
}

func (r *DOMRenderer) NewComponent(name string) render.Component {
	return &WebComponent{
		name: name,
	}
}

func (r *DOMRenderer) Render(comps ...render.Component) {
	for _, c := range comps {
		webC, ok := c.(*WebComponent)
		if !ok {
			panic("invalid comp type")
		}

		frag, _ := r.createElement(webC)
		r.body.Call("appendChild", frag)
	}
}

func (r *DOMRenderer) createElement(c *WebComponent) (frag, el js.Value) {
	frag = r.doc.Call("createDocumentFragment")

	el = r.doc.Call("createElement", c.name)
	c.el = &el

	if c.text != "" {
		el.Set("innerText", c.text)
	}

	for key, val := range c.attrs {
		if strings.HasPrefix(key, "on:") {
			valFunc, ok := val.(func())
			if !ok {
				log.Fatal("only func allowed for eventlisteners")
			}

			jsFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
				valFunc()
				return nil
			})
			// todo: handle release
			el.Call("addEventListener", strings.TrimPrefix(key, "on:"), jsFunc)

			continue
		}

		if strings.HasPrefix(key, "bind:") {
			valFunc, ok := val.(func(val any))
			if !ok {
				log.Fatal("expected func(val any) when binding")
			}

			attr := strings.TrimPrefix(key, "bind:")
			switch attr {
			case "value":
				jsFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
					valFunc(el.Get("value").String())
					return nil
				})
				// todo: handle release
				el.Call("addEventListener", "input", jsFunc)
			default:
				panic("bind to " + attr + "not implemented")
			}
			continue
		}

		c.SetAttribute(key, val)
	}

	for _, cc := range c.comps {
		webCC, ok := cc.(*WebComponent)
		if !ok {
			panic("invalid comp type")
		}

		cfrag, cel := r.createElement(webCC)
		webCC.el = &cel
		el.Call("appendChild", cfrag)
	}

	frag.Call("appendChild", el)
	return frag, el
}

func (r *DOMRenderer) Append(parent render.Component, child render.Component) {
	p, ok := parent.(*WebComponent)
	if !ok {
		panic("invalid parent comp type")
	}

	c, ok := child.(*WebComponent)
	if !ok {
		panic("invalid parent comp type")
	}

	if p.el == nil {
		panic("nil parent el")
	}

	if c.el != nil {
		p.el.Call("appendChild", *c.el)
		return
	}

	frag, el := r.createElement(c)
	c.el = &el
	p.el.Call("appendChild", frag)
}
