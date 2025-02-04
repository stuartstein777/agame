package main

import (
	"bufio"
	"os"

	"github.com/stuartstein777/roguelike/entities"
)

const (
	Open = iota
	Wall
	Blank
	ClosedDoor
	OpenDoor
)

func LoadDungeon(dungeon [][]entities.Tile) {
	file := "dungeon1.txt"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, char := range line {
			switch char {
			case '.':
				dungeon[y][x].TileType = Open
			case '@':
				dungeon[y][x].TileType = Wall
			case ' ':
				dungeon[y][x].TileType = Blank
			case 'K':
				dungeon[y][x].TileType = Open
				dungeon[y][x].Items = append(dungeon[y][x].Items, entities.Item{Type: entities.Key, SpritesheetIdx: 6})
			case 'D':
				dungeon[y][x].TileType = ClosedDoor
			}
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func MakeDungeon() [][]entities.Tile {
	dungeon := make([][]entities.Tile, 100)

	for y := 0; y < 100; y++ {
		dungeon[y] = make([]entities.Tile, 100)

		for x := 0; x < 100; x++ {
			dungeon[y][x].TileType = Open
			dungeon[y][x].X = x
			dungeon[y][x].Y = y
		}
	}

	LoadDungeon(dungeon)
	return dungeon
}

func MakeLayers(dungeon [][]entities.Tile) [][]int {
	layers := make([][]int, 3)

	for i := 0; i < 3; i++ {
		layers[i] = make([]int, 100*100)
	}

	// set all values in layers to 0
	for i := 0; i < 100*100; i++ {
		layers[0][i] = 0
		layers[1][i] = 0
	}

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			layers[1][x+(y*100)] = 5
			layers[2][x+(y*100)] = 5
		}
	}

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if dungeon[y][x].TileType == Blank {
				layers[0][x+(y*100)] = 5
			} else if dungeon[y][x].TileType == Open {
				layers[0][x+(y*100)] = 1
				if len(dungeon[y][x].Items) != 0 {
					item := dungeon[y][x].Items[0]
					layers[1][x+(y*100)] = item.SpritesheetIdx
				}
			} else if dungeon[y][x].TileType == ClosedDoor {
				layers[0][x+(y*100)] = 1
				layers[1][x+(y*100)] = 7
			} else {
				if y < 99 && dungeon[y+1][x].TileType == Wall {
					layers[0][x+(y*100)] = 0
				} else {
					layers[0][x+(y*100)] = 3
				}
			}
		}
	}

	return layers
}
