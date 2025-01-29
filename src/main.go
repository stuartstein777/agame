package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 704
	screenHeight = 480
)

const (
	tileSize = 32
)

var (
	tilesImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
}

type Coord struct {
	X int
	Y int
}

type ViewPort struct {
	xy     Coord
	height int
	width  int
}
type Game struct {
	playLocation Coord
	viewPort     ViewPort
	layers       [][]int
	dungeon      [][]Tile
}

func (g *Game) Update() error {

	location := ((g.playLocation.Y * 100) - 1) + g.playLocation.X //todo: is this wrong?

	g.layers[1][location] = 5 // 5 = transparent image.

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.playLocation.X -= 1
		if g.playLocation.X == 0 {
			g.playLocation.X = 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.playLocation.X += 1
		if g.playLocation.X == 23 {
			g.playLocation.X = 22
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.playLocation.Y += 1
		if g.playLocation.Y == 15 {
			g.playLocation.Y = 14
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.playLocation.Y -= 1
		if g.playLocation.Y < 0 {
			g.playLocation.Y = 0
		}
	}

	// update player spritee in layers [1]
	location = ((g.playLocation.Y * 22) - 1) + g.playLocation.X
	g.layers[1][location] = 2

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w := tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	for _, l := range g.layers {

		for y := 0; y < g.viewPort.height; y++ {
			viewPortStart := (((g.viewPort.xy.Y + y) * 100) - 1) + g.viewPort.xy.X
			viewPortEnd := viewPortStart + g.viewPort.width

			row := l[viewPortStart:viewPortEnd]

			for i, t := range row {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64((i * tileSize)), float64(y*tileSize))

				sx := (t % tileXCount) * tileSize
				sy := (t / tileXCount) * tileSize
				screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
			}
		}
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("Player location: %v", g.playLocation))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	dungeon := MakeDungeon()
	layers := MakeLayers(dungeon)

	g := &Game{
		playLocation: Coord{
			X: 30,
			Y: 19,
		},
		viewPort: ViewPort{
			xy: Coord{
				X: 20,
				Y: 12,
			},
			height: 15,
			width:  22,
		},
		dungeon: MakeDungeon(),
		layers:  layers,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rogue like")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
