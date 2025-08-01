package ast

import (
	source "go/token"

	"github.com/tifye/flamingo/assert"
)

type NodeType string

const (
	NodeTypeFrag NodeType = "fragment"
)

type Node interface {
	Pos() source.Pos
}

type RenderNode interface {
	Node
	elementNode()
}

type (
	// A File node represents a single Flamingo component
	File struct {
		CodeBlock *CodeBlock
		Fragment  *Fragment
	}

	CodeBlock struct {
		TopFence    source.Pos
		BottomFence source.Pos
	}

	Fragment struct {
		Nodes []RenderNode
	}

	Element struct {
		LeftChevron  source.Pos // left chevron of opening tag
		RightChevron source.Pos // right chevron of close tag
		Name         *Ident
		Attrs        []*Attribute
		Nodes        []RenderNode
	}

	Ident struct {
		Position source.Pos // start of the name
		Name     string
	}

	Attribute struct {
		Name         *Ident
		ValueLiteral string
	}

	Text struct {
		Position source.Pos
		Literal  string
	}
)

func (n *File) Pos() source.Pos {
	if n.CodeBlock != nil {
		return n.CodeBlock.TopFence
	}
	if n.Fragment != nil {
		return n.Fragment.Pos()
	}
	panic("node has no source location")
}
func (n *CodeBlock) Pos() source.Pos { return n.TopFence }
func (n *Fragment) Pos() source.Pos {
	l := len(n.Nodes)
	assert.Assert(l > 0, "fragment has no children")
	return n.Nodes[l-1].Pos()
}
func (n *Element) Pos() source.Pos   { return n.LeftChevron }
func (n *Ident) Pos() source.Pos     { return n.Position }
func (n *Attribute) Pos() source.Pos { return n.Name.Pos() }
func (n *Text) Pos() source.Pos      { return n.Position }

// elementNode() makes sure that only element nodes can be assigned to an Element
func (*Element) elementNode()  {}
func (*Text) elementNode()     {}
func (*Fragment) elementNode() {}
