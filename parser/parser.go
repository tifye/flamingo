package parser

import (
	"fmt"
	"strings"

	"github.com/tifye/flamingo/assert"
	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
	"github.com/tifye/flamingo/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken() // sets peekToken
	p.nextToken() // sets curToken
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() *ast.File {
	root := &ast.File{}
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
		return &ast.Text{Lit: strings.TrimSpace(p.curToken.Literal)}
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
		Attrs: make([]*ast.Attr, 0),
		Nodes: make([]ast.Element, 0),
	}

	for p.tryPeek(token.IDENT) {
		attr := p.parseAttribute()
		if attr != nil {
			comp.Attrs = append(comp.Attrs, attr)
		}
	}

	assert.Assert(p.isPeekToken(token.RIGHT_CHEV), "expected next token to be a right chevron")
	_ = p.expectPeek(token.RIGHT_CHEV)

	for {
		p.nextToken()
		el := p.parseElement()
		if el == nil {
			break
		}
		comp.Nodes = append(comp.Nodes, el)
	}

	if !p.expectPeek(token.SLASH) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	if p.curToken.Literal != comp.Name.Name {
		p.errorf("unexpected closing tag %s, expected %s", p.curToken.Literal, comp.Name.Name)
		return nil
	}

	if !p.expectPeek(token.RIGHT_CHEV) {
		return nil
	}

	return comp
}

func (p *Parser) parseAttribute() *ast.Attr {
	assert.Assert(p.isCurToken(token.IDENT), "expected curToken to be an identifier")

	attr := &ast.Attr{
		Name: &ast.Ident{Name: p.curToken.Literal},
	}

	if !p.tryPeek(token.ASSIGN) {
		attr.ValueLit = "true"
		return attr
	}

	if !p.expectPeek(token.QUOTE) {
		return nil
	}

	if p.tryPeek(token.TEXT) {
		attr.ValueLit = p.curToken.Literal
	}

	if !p.expectPeek(token.QUOTE) {
		return nil
	}

	return attr
}

func (p *Parser) isCurToken(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// tryPeek advances the Parser ahead by one token
// if p.peekToken equals the passed TokenType.
func (p *Parser) tryPeek(t token.TokenType) bool {
	if p.isPeekToken(t) {
		p.nextToken()
		return true
	}
	return false
}

// expectPeek behaves similar to tryPeek but generates an error
// if p.peekToken is not equal to the passed TokenType.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if !p.tryPeek(t) {
		p.peekError(t)
		return false
	}
	return true
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) errorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	p.errors = append(p.errors, msg)
}
