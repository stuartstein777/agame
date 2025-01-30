package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/stuartstein777/roguelike/entities"
	"github.com/stuartstein777/roguelike/resources"
)

const (
	screenWidth  = 704
	screenHeight = 480
)

const (
	tileSize = 32
)

type Game struct {
	playLocation entities.Coord
	viewPort     entities.ViewPort
	layers       [][]int
	dungeon      [][]entities.Tile
	debugMessage string
}

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(resources.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	resources.TilesImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Update() error {

	location := (g.playLocation.Y * 100) + g.playLocation.X

	g.layers[1][location] = 5 // 5 = transparent image.

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.debugMessage = "Left"
		if g.dungeon[g.playLocation.Y][g.playLocation.X-1].TileType != Wall {
			g.playLocation.X -= 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.debugMessage = "Up"
		if g.dungeon[g.playLocation.Y-1][g.playLocation.X].TileType != Wall {
			g.playLocation.Y -= 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.debugMessage = "Down"

		if g.dungeon[g.playLocation.Y+1][g.playLocation.X].TileType != Wall {
			g.playLocation.Y += 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.debugMessage = "Right"
		if g.dungeon[g.playLocation.Y][g.playLocation.X+1].TileType != Wall {
			g.playLocation.X += 1
		}
	}

	// update player sprite in layers [1]
	location = (g.playLocation.Y * 100) + g.playLocation.X
	g.layers[1][location] = 2

	// if player is within 3 tiles of the edge of the viewport, move the viewport
	// in the direction the player is moving.
	if g.playLocation.X-g.viewPort.XY.X < 3 {
		g.viewPort.XY.X -= 1
	} else if g.playLocation.X > g.viewPort.XY.X+g.viewPort.Width-5 { //BUG: why am i setting this to -5 rather than -3?
		g.viewPort.XY.X += 1
	} else if g.playLocation.Y > g.viewPort.XY.Y+g.viewPort.Height-3 {
		g.viewPort.XY.Y += 1
	} else if g.playLocation.Y-g.viewPort.XY.Y < 3 {
		g.viewPort.XY.Y -= 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w := resources.TilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	for _, l := range g.layers {

		for y := 0; y < g.viewPort.Height; y++ {
			viewPortStart := (((g.viewPort.XY.Y + y) * 100) - 1) + g.viewPort.XY.X
			viewPortEnd := viewPortStart + g.viewPort.Width

			row := l[viewPortStart:viewPortEnd]

			for i, t := range row {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64((i * tileSize)), float64(y*tileSize))

				sx := (t % tileXCount) * tileSize
				sy := (t / tileXCount) * tileSize
				screen.DrawImage(resources.TilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
			}
		}
	}

	gridImage := ebiten.NewImage(g.viewPort.Width*tileSize, g.viewPort.Height*tileSize)
	gridColor := color.RGBA{R: 0, G: 0, B: 200, A: 255} // Light gray grid

	for y := 0; y <= g.viewPort.Height; y++ {
		vector.StrokeLine(gridImage, 0, float32(y*tileSize), float32(g.viewPort.Width*tileSize), float32(y*tileSize), 0.5, gridColor, true)
	}
	for x := 0; x <= g.viewPort.Width; x++ {
		vector.StrokeLine(gridImage, float32(x*tileSize), 0, float32(x*tileSize), float32(g.viewPort.Height*tileSize), 0.5, gridColor, true)
	}

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(gridImage, op)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 5, 5)
	ebitenutil.DebugPrintAt(screen,
		fmt.Sprintf("Player location: %v, ViewPort location :%v, Pushed: %s", g.playLocation, g.viewPort, g.debugMessage),
		5, 20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	dungeon := MakeDungeon()
	layers := MakeLayers(dungeon)

	g := &Game{
		playLocation: entities.Coord{
			X: 30,
			Y: 19,
		},
		viewPort: entities.ViewPort{
			XY: entities.Coord{
				X: 20,
				Y: 15,
			},
			Height: 15,
			Width:  22,
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
