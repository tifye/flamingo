package meep

import (
	"fmt"

	"github.com/tifye/flamingo/render"
)

var (
	izu = "izu"
)

// prefixed with component name
type Mino struct {
	Meep string `prop`

	text string
}

func (c *Mino) OnMount() {
	meep()
}

func meep() {
	fmt.Println("mino")
}

type MinoProps struct {
	Meep string
}

func MinoComp(renderer render.Renderer, props MinoProps) {
	comp := &Mino{
		Meep: props.Meep,
	}

	// ... or other lifecycle stuff
	if mounter, ok := comp.(Mounter); ok {
		mounter.OnMount()
	}

	// ...
}
