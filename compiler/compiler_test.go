package compiler

import (
	"fmt"
	gtoken "go/token"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	fileb, err := os.ReadFile("testdata/Mino.flamingo")
	require.NoError(t, err)

	fset := gtoken.NewFileSet()
	file := fset.AddFile("Mino", fset.Base(), len(fileb))
	output := &strings.Builder{}
	err = CompileFile("main", file, string(fileb), output)
	assert.NoError(t, err)
	fmt.Println(output.String())
}
