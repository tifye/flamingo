package ast

type NodeType string

const (
	NodeTypeFrag NodeType = "fragment"
)

type Node interface {
	TokenLit() string
	// Pos() token.Pos
}

type Element interface {
	Node
	elementNode()
}

type (
	// A File node represents a single Flamingo component
	File struct {
		Fragment *Fragment
	}

	Fragment struct {
		Nodes []Element
	}

	Component struct {
		Name  *Ident
		Attrs []*Attr
		Nodes []Element
	}

	Ident struct {
		Name string
	}

	Attr struct {
		Name     *Ident
		ValueLit string
	}

	Text struct {
		Lit string
	}
)

func (c *Component) TokenLit() string { return c.Name.Name }
func (i *Ident) TokenLit() string     { return i.Name }
func (t *Text) TokenLit() string      { return t.Lit }
func (a *Attr) TokenLit() string      { return a.Name.Name + " " + a.ValueLit }
func (r *Fragment) TokenLit() string {
	if len(r.Nodes) > 0 {
		return r.Nodes[0].TokenLit()
	}
	return ""
}
func (r *File) TokenLit() string {
	if r.Fragment != nil {
		return r.Fragment.TokenLit()
	}
	return ""
}

// elementNode() makes sure that only element nodes can be assigned to an Element
func (*Component) elementNode() {}
func (*Text) elementNode()      {}
func (*Fragment) elementNode()  {}
