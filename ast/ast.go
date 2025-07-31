package ast

type NodeType string

const (
	NodeTypeFrag NodeType = "fragment"
)

type Node interface {
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

// elementNode() makes sure that only element nodes can be assigned to an Element
func (*Component) elementNode() {}
func (*Text) elementNode()      {}
func (*Fragment) elementNode()  {}
