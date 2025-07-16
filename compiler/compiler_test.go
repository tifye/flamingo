package compiler

import (
	"fmt"
	gtoken "go/token"
	"testing"
)

func TestCompiler(t *testing.T) {
	input := `<div class="bg-rose-500"><span>mino</span></div>`
	fset := gtoken.NewFileSet()
	file := fset.AddFile("Mino", fset.Base(), len(input))
	output := Compile("test", file, input)
	fmt.Println(output)
}
