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
	player       entities.Player
	viewPort     entities.ViewPort
	layers       [][]int
	dungeon      [][]entities.Tile
	debugMessage string
	showGrid     bool
	enemies      []entities.Enemy
}

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(resources.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	resources.TilesImage = ebiten.NewImageFromImage(img)
}

func PlayerHasKey(g *Game) bool {
	for _, item := range g.player.Items {
		if item.Type == entities.Key {
			return true
		}
	}
	return false
}

func CanPlayerMove(g *Game, x, y int) bool {
	if g.dungeon[y][x].TileType == Wall {
		return false
	}
	if g.dungeon[y][x].TileType == ClosedDoor {
		if !PlayerHasKey(g) {
			return false
		} else {
			// remove key from inventory
			for i, item := range g.player.Items {
				if item.Type == entities.Key {
					g.player.Items = append(g.player.Items[:i], g.player.Items[i+1:]...)
				}
			}

			// update the dungeon CloseDoor to OpenDoor
			// update the layer to the open door sprite
			g.dungeon[y][x].TileType = OpenDoor
		}
	}
	return true
}

func (g *Game) Update() error {

	location := (g.player.PlayLocation.Y * 100) + g.player.PlayLocation.X

	g.layers[2][location] = 5

	if len(g.dungeon[g.player.PlayLocation.Y][g.player.PlayLocation.X].Items) != 0 {
		item := g.dungeon[g.player.PlayLocation.Y][g.player.PlayLocation.X].Items[0]
		g.player.Items = append(g.player.Items, item)
		g.dungeon[g.player.PlayLocation.Y][g.player.PlayLocation.X].Items = nil
		g.layers[1][location] = 5
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.debugMessage = "Left"
		if CanPlayerMove(g, g.player.PlayLocation.X-1, g.player.PlayLocation.Y) {
			g.player.PlayLocation.X -= 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.debugMessage = "Up"
		if CanPlayerMove(g, g.player.PlayLocation.X, g.player.PlayLocation.Y-1) {
			g.player.PlayLocation.Y -= 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.debugMessage = "Down"
		if CanPlayerMove(g, g.player.PlayLocation.X, g.player.PlayLocation.Y+1) {
			g.player.PlayLocation.Y += 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.debugMessage = "Right"
		if CanPlayerMove(g, g.player.PlayLocation.X+1, g.player.PlayLocation.Y) {
			g.player.PlayLocation.X += 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.showGrid = !g.showGrid
	}

	// update player sprite in layers [1]
	location = (g.player.PlayLocation.Y * 100) + g.player.PlayLocation.X
	g.layers[2][location] = 2

	// if player is within 3 tiles of the edge of the viewport, move the viewport
	// in the direction the player is moving.
	if g.player.PlayLocation.X-g.viewPort.XY.X < 3 {
		g.viewPort.XY.X -= 1
	} else if g.player.PlayLocation.X > g.viewPort.XY.X+g.viewPort.Width-5 { //BUG: why am i setting this to -5 rather than -3?
		g.viewPort.XY.X += 1
	} else if g.player.PlayLocation.Y > g.viewPort.XY.Y+g.viewPort.Height-3 {
		g.viewPort.XY.Y += 1
	} else if g.player.PlayLocation.Y-g.viewPort.XY.Y < 3 {
		g.viewPort.XY.Y -= 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w := resources.TilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	for lidx, l := range g.layers {

		for y := 0; y < g.viewPort.Height; y++ {
			viewPortStart := ((g.viewPort.XY.Y + y) * 100) + g.viewPort.XY.X
			viewPortEnd := viewPortStart + g.viewPort.Width

			row := l[viewPortStart:viewPortEnd]

			for i, t := range row {

				sx := (t % tileXCount) * tileSize
				sy := (t / tileXCount) * tileSize

				dungeonX := g.viewPort.XY.X + i
				dungeonY := g.viewPort.XY.Y + y

				currentTile := g.dungeon[dungeonY][dungeonX]

				if lidx == 1 && t == 7 {

					// Calculate dungeon coordinates

					// Check for walls above and below the door

					// check if the dungeon tile is a closed door or an open door
					if currentTile.TileType == ClosedDoor {
						if g.dungeon[dungeonY-1][dungeonX].TileType == Wall &&
							g.dungeon[dungeonY+1][dungeonX].TileType == Wall {
							op := &ebiten.DrawImageOptions{}

							op.GeoM.Translate(float64(tileSize)/2, float64(tileSize)/2)

							// Rotate 90 degrees clockwise (swap x and y)
							op.GeoM.SetElement(0, 0, 0)
							op.GeoM.SetElement(0, 1, -1)
							op.GeoM.SetElement(1, 0, 1)
							op.GeoM.SetElement(1, 1, 0)

							// Move back to tile's correct position on screen
							op.GeoM.Translate(float64(i*tileSize)+float64(tileSize)/2, float64(y*tileSize)-float64(tileSize)/2)
							screen.DrawImage(resources.TilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
						} else {
							op := &ebiten.DrawImageOptions{}
							op.GeoM.Translate(float64(i*tileSize), float64(y*tileSize))
							screen.DrawImage(resources.TilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
						}
					} else if currentTile.TileType == OpenDoor {
						if g.dungeon[dungeonY][dungeonX-1].TileType == Wall &&
							g.dungeon[dungeonY][dungeonX+1].TileType == Wall {
							op := &ebiten.DrawImageOptions{}

							// Move to the hinge point (top-center of the tile)
							op.GeoM.Translate(float64(tileSize)/2, float64(tileSize)/2)

							// Rotate 90 degrees clockwise (swap x and y)
							op.GeoM.SetElement(0, 0, 0)
							op.GeoM.SetElement(0, 1, -1)
							op.GeoM.SetElement(1, 0, 1)
							op.GeoM.SetElement(1, 1, 0)

							// Move back to tile's correct position on screen
							op.GeoM.Translate(float64(i*tileSize)+float64(tileSize)/2, float64(y*tileSize)-float64(tileSize)/2)

							// Draw the door
							screen.DrawImage(resources.TilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
						} else {
							op := &ebiten.DrawImageOptions{}
							op.GeoM.Translate(float64(i*tileSize), float64(y*tileSize))
							screen.DrawImage(resources.TilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
						}
					}
				} else {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(i*tileSize), float64(y*tileSize))
					screen.DrawImage(resources.TilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
				}
			}
		}
	}

	if g.showGrid {
		gridImage := ebiten.NewImage(g.viewPort.Width*tileSize, g.viewPort.Height*tileSize)
		gridColor := color.RGBA{R: 48, G: 213, B: 200, A: 255} // Light gray grid

		for y := 0; y <= g.viewPort.Height; y++ {
			yPos := y * tileSize
			vector.StrokeLine(gridImage, 0, float32(y*tileSize), float32(g.viewPort.Width*tileSize), float32(y*tileSize), 0.5, gridColor, true)
			if yPos > 0 {
				ebitenutil.DebugPrintAt(gridImage, fmt.Sprintf("%d", g.viewPort.XY.Y+y), 5, yPos)
			}
		}
		for x := 0; x <= g.viewPort.Width; x++ {
			xPos := x * tileSize
			vector.StrokeLine(gridImage, float32(x*tileSize), 0, float32(x*tileSize), float32(g.viewPort.Height*tileSize), 0.5, gridColor, true)
			if x < g.viewPort.Width && xPos > 0 {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", g.viewPort.XY.X+x-1), xPos+5, 2)
			}
		}

		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(gridImage, op)

		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("Player location: %v, ViewPort location :%v, Pushed: %s", g.player.PlayLocation, g.viewPort, g.debugMessage),
			60, 50)

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Player items: %v", g.player.Items), 60, 30)
		//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 5, 5)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	dungeon := MakeDungeon()
	layers := MakeLayers(dungeon)

	g := &Game{
		player: entities.Player{
			PlayLocation: entities.Coord{
				X: 30,
				Y: 19,
			},
		},
		viewPort: entities.ViewPort{
			XY: entities.Coord{
				X: 19,
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
