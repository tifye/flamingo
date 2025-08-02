package compiler

import (
	"fmt"
	gtoken "go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tifye/flamingo/parser"
)

func TestCompiler(t *testing.T) {
	fset := gtoken.NewFileSet()
	output := &strings.Builder{}
	root, err := parser.ParseFile(fset, "testdata/Mino.flamingo", nil)
	assert.NoError(t, err)

	err = CompileFile("main", "Mino", root, output)
	assert.NoError(t, err)

	fmt.Println(output.String())
}
