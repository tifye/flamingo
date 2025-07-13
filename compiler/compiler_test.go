package compiler

import (
	"fmt"
	"testing"
)

func TestCompiler(t *testing.T) {
	input := `<div class="bg-rose-500"><span>mino</span></div>`
	output := Compile(input)
	fmt.Println(output)
}
