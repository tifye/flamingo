package parser

import (
	"bytes"
	"errors"
	"fmt"
	source "go/token"
	"io"
	"os"
	"strings"

	"github.com/tifye/flamingo/assert"
	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
	"github.com/tifye/flamingo/token"
)

func ParseElement(src any) (*ast.Element, error) {
	input, err := readSource("", src)
	if err != nil {
		return nil, err
	}

	fset := source.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(input))
	l := lexer.NewLexer(file, string(input)).WithState(lexer.LexTag)
	p := NewParser(l)
	el := p.parseElement()
	if n := len(p.Errors()); n > 0 {
		return el, errors.New(strings.Join(p.Errors(), "; "))
	}
	return el, nil
}

func ParseFile(fset *source.FileSet, filename string, src any) (*ast.File, error) {
	input, err := readSource(filename, src)
	if err != nil {
		return nil, fmt.Errorf("reading source: %s", err)
	}

	file := fset.AddFile(filename, fset.Base(), len(input))
	l := lexer.NewLexer(file, string(input))
	p := NewParser(l)

	fileNode := p.Parse()
	if len(p.errors) == 0 {
		return fileNode, nil
	}
	return fileNode, fmt.Errorf("%v", p.errors)
}

func readSource(filename string, src any) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return io.ReadAll(s)
		}

		return nil, errors.New("invalid source type")
	}

	return os.ReadFile(filename)
}

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
		Nodes: make([]ast.RenderNode, 0),
	}

	for !p.isCurToken(token.EOF) {
		el := p.parseRenderNode()
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
func (p *Parser) parseRenderNode() ast.RenderNode {
	switch p.curToken.Type {
	case token.TEXT:
		return &ast.Text{
			Position: p.curToken.Pos,
			Literal:  strings.TrimSpace(p.curToken.Literal),
		}
	case token.LEFT_CHEVRON:
		if p.isPeekToken(token.SLASH) {
			return nil
		}
		return p.parseElement()
	default:
		return nil
	}
}

func (p *Parser) parseElement() (el *ast.Element) {
	assert.Assert(p.isCurToken(token.LEFT_CHEVRON), "expected left chevron")
	assert.Assert(p.isPeekToken(token.IDENT), "expected next token to be an identifier")

	if !p.tryPeek(token.IDENT) {
		return nil
	}

	comp := &ast.Element{
		LeftChevron: p.curToken.Pos - 1,
		Name: &ast.Ident{
			Position: p.curToken.Pos,
			Name:     p.curToken.Literal,
		},
		Attrs: make([]*ast.Attribute, 0),
		Nodes: make([]ast.RenderNode, 0),
	}
	defer func() {
		if el != nil {
			assert.Assert(comp.RightChevron.IsValid(), "expected right chevron location to be set")
		}
	}()

	for p.tryPeek(token.IDENT) {
		attr := p.parseAttribute()
		if attr != nil {
			comp.Attrs = append(comp.Attrs, attr)
		}
	}

	assert.Assert(p.isPeekToken(token.RIGHT_CHEVRON), "expected next token to be a right chevron")
	_ = p.expectPeek(token.RIGHT_CHEVRON)

	for {
		p.nextToken()
		el := p.parseRenderNode()
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

	if !p.expectPeek(token.RIGHT_CHEVRON) {
		return nil
	}

	comp.RightChevron = p.curToken.Pos

	return comp
}

func (p *Parser) parseAttribute() *ast.Attribute {
	assert.Assert(p.isCurToken(token.IDENT), "expected curToken to be an identifier")

	attr := &ast.Attribute{
		Name: &ast.Ident{
			Position: p.curToken.Pos,
			Name:     p.curToken.Literal,
		},
	}

	if !p.tryPeek(token.ASSIGN) {
		attr.ValueLiteral = "true"
		return attr
	}

	if !p.expectPeek(token.QUOTE) {
		return nil
	}

	if p.tryPeek(token.TEXT) {
		attr.ValueLiteral = p.curToken.Literal
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
// If successful p.curToken is guaranteed to contain
// a token of type = t.
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
