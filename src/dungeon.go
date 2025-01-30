package main

import (
	"bufio"
	"os"

	"github.com/stuartstein777/roguelike/entities"
)

const (
	Open = iota
	Wall = iota
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
			case '#':
				dungeon[y][x].TileType = Open
			case '@':
				dungeon[y][x].TileType = Wall
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
	layers := make([][]int, 2)

	for i := 0; i < 2; i++ {
		layers[i] = make([]int, 100*100)
	}

	// set all values in layers to 0
	for i := 0; i < 100*100; i++ {
		layers[0][i] = 0
		layers[1][i] = 0
	}

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if dungeon[y][x].TileType == Open {
				layers[0][x+(y*100)] = 1
			} else {
				layers[0][x+(y*100)] = 0
			}
		}
	}

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			layers[1][x+(y*100)] = 5
		}
	}

	return layers
}
