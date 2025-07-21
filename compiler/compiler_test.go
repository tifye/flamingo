package compiler

import (
	"fmt"
	gtoken "go/token"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tifye/flamingo/lexer"
	"github.com/tifye/flamingo/parser"
)

func TestCompiler(t *testing.T) {
	fileb, err := os.ReadFile("testdata/Mino.flamingo")
	require.NoError(t, err)

	fset := gtoken.NewFileSet()
	file := fset.AddFile("Mino", fset.Base(), len(fileb))
	output := &strings.Builder{}

	l := lexer.NewLexer(file, string(fileb))
	p := parser.NewParser(l)
	root := p.Parse()

	assert.Empty(t, p.Errors())

	err = CompileFile("main", "Mino", root, string(fileb), output)
	assert.NoError(t, err)
	fmt.Println(output.String())
}
