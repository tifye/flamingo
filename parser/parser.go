package parser

import (
	"github.com/tifye/flamingo/assert"
	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
	"github.com/tifye/flamingo/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken() // sets peekToken
	p.nextToken() // sets curToken
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() *ast.Root {
	root := &ast.Root{}
	root.Fragment = &ast.Fragment{
		Nodes: make([]ast.Element, 0),
	}

	for !p.isCurToken(token.EOF) {
		el := p.parseElement()
		if el != nil {
			root.Fragment.Nodes = append(root.Fragment.Nodes, el)
		}

		if p.isPeekToken(token.EOF) {
			break
		}

		p.nextToken()
	}

	return root
}

// Tries to parse curToken to an Element
// otherwise returns nil.
//
// Uses peekToken incases of '<' specifically
// looking for '/' in which case it also returns nil.
func (p *Parser) parseElement() ast.Element {
	switch p.curToken.Type {
	case token.TEXT:
		return &ast.Text{Lit: p.curToken.Literal}
	case token.LEFT_CHEV:
		if p.isPeekToken(token.SLASH) {
			return nil
		}
		return p.parseComponent()
	default:
		return nil
	}
}

func (p *Parser) parseComponent() *ast.Component {
	assert.Assert(p.isCurToken(token.LEFT_CHEV), "expected left chevron")
	assert.Assert(p.isPeekToken(token.IDENT), "expected next token to be an identifier")

	if !p.tryPeek(token.IDENT) {
		return nil
	}

	comp := &ast.Component{
		Name:  &ast.Ident{Name: p.curToken.Literal},
		Nodes: make([]ast.Element, 0),
	}

	if p.isPeekToken(token.IDENT) {
		// parse attrs
	}

	assert.Assert(p.isPeekToken(token.RIGHT_CHEV), "expected next token to be a right chevron")
	_ = p.tryPeek(token.RIGHT_CHEV)

	for {
		p.nextToken()
		el := p.parseElement()
		if el == nil {
			break
		}
		comp.Nodes = append(comp.Nodes, el)
	}

	if !p.tryPeek(token.SLASH) {
		panic("expected slash")
	}

	if !p.tryPeek(token.IDENT) {
		panic("expected identifier")
	}

	if p.curToken.Literal != comp.Name.Name {
		panic("invalid closing tag")
	}

	if !p.tryPeek(token.RIGHT_CHEV) {
		panic("expected right chevron")
	}

	return comp
}

func (p *Parser) isCurToken(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) tryPeek(t token.TokenType) bool {
	if p.isPeekToken(t) {
		p.nextToken()
		return true
	}

	return false
}
