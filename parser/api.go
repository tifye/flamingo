package parser

import (
	"bytes"
	"errors"
	"fmt"
	gtoken "go/token"
	"io"
	"os"

	"github.com/tifye/flamingo/ast"
	"github.com/tifye/flamingo/lexer"
)

func ParseFile(fset *gtoken.FileSet, filename string, src any) (*ast.File, error) {
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
