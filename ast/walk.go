package ast

import "reflect"

type Visitor interface {
	Visit(node Node) (w Visitor)
}

func walkList[N Node](v Visitor, list []N) {
	for _, node := range list {
		Walk(v, node)
	}
}

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *File:
		Walk(v, n.Fragment)
	case *Fragment:
		walkList(v, n.Nodes)
	case *Element:
		Walk(v, n.Name)
		walkList(v, n.Attrs)
		walkList(v, n.Nodes)
	case *Attribute:
		Walk(v, n.Name)
	case *Text:
	case *Ident:
	default:
		panic("cannot walk node of type: " + reflect.TypeOf(n).String())
	}
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
