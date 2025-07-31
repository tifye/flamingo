package ast

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
	case *Tag:
		Walk(v, n.Name)
		walkList(v, n.Attrs)
		walkList(v, n.Nodes)
	case *Attr:
		Walk(v, n.Name)
	case *Text:
	case *Ident:
	default:
		return
	}
}
