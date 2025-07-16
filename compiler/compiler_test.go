package compiler

import (
	"fmt"
	gtoken "go/token"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	fileb, err := os.ReadFile("testdata/Mino.flamingo")
	require.NoError(t, err)

	fset := gtoken.NewFileSet()
	file := fset.AddFile("Mino", fset.Base(), len(fileb))
	output := Compile("test", file, string(fileb))
	fmt.Println(output)
}
