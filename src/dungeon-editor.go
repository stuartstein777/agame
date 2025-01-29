package main

// import (
// 	"fmt"
// 	"image/color"
// 	"os"
// 	"strings"

// 	"github.com/hajimehoshi/ebiten/v2"
// )

// const (
// 	gridWidth  = 100
// 	gridHeight = 100
// 	tileSize   = 32
// )

// type Editor struct {
// 	grid [gridHeight][gridWidth]rune // '#' for walls, ' ' for empty
// }

// func (e *Editor) Update() error {
// 	// Get mouse position
// 	x, y := ebiten.CursorPosition()
// 	gridX := x / tileSize
// 	gridY := y / tileSize

// 	// Toggle grid on left mouse click
// 	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
// 		if gridX >= 0 && gridX < gridWidth && gridY >= 0 && gridY < gridHeight {
// 			e.grid[gridY][gridX] = '#'
// 		}
// 	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
// 		if gridX >= 0 && gridX < gridWidth && gridY >= 0 && gridY < gridHeight {
// 			e.grid[gridY][gridX] = ' '
// 		}
// 	}

// 	// Save the map
// 	if ebiten.IsKeyPressed(ebiten.KeyS) {
// 		e.saveMap("dungeon.txt")
// 	}

// 	// Load the map
// 	if ebiten.IsKeyPressed(ebiten.KeyL) {
// 		e.loadMap("dungeon.txt")
// 	}

// 	return nil
// }

// func (e *Editor) Draw(screen *ebiten.Image) {
// 	for y := 0; y < gridHeight; y++ {
// 		for x := 0; x < gridWidth; x++ {
// 			col := ebiten.ColorM{}
// 			if e.grid[y][x] == '#' {
// 				col.Scale(1, 1, 0, 1) // Green for walls
// 			} else {
// 				col.Scale(1, 1, 1, 1) // White for empty
// 			}
// 			rect := ebiten.NewImage(tileSize, tileSize)
// 			rect.Fill(col.Apply(color.RGBA{R: 255, G: 255, B: 255, A: 100}))
// 			op := &ebiten.DrawImageOptions{}
// 			op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
// 			screen.DrawImage(rect, op)
// 		}
// 	}
// }

// func (e *Editor) saveMap(filename string) {
// 	var sb strings.Builder
// 	for _, row := range e.grid {
// 		sb.WriteString(string(row[:]) + "\n")
// 	}
// 	err := os.WriteFile(filename, []byte(sb.String()), 0644)
// 	if err != nil {
// 		fmt.Println("Error saving:", err)
// 	}
// }

// func (e *Editor) loadMap(filename string) {
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		fmt.Println("Error loading:", err)
// 		return
// 	}
// 	lines := strings.Split(string(data), "\n")
// 	for y, line := range lines {
// 		if y >= gridHeight {
// 			break
// 		}
// 		copy(e.grid[y][:], []rune(line))
// 	}
// }

// func (e *Editor) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return gridWidth * tileSize, gridHeight * tileSize
// }

// func main() {
// 	ebiten.SetWindowSize(gridWidth*tileSize, gridHeight*tileSize)
// 	ebiten.SetWindowTitle("Dungeon Editor")
// 	editor := &Editor{}
// 	if err := ebiten.RunGame(editor); err != nil {
// 		panic(err)
// 	}
// }
