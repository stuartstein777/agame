package resources

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	TilesImage *ebiten.Image
	//go:embed art/spritesheet.png
	Tiles_png []byte
)
