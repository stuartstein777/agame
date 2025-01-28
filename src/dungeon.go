package main

const (
	Open  = iota
	Floor = iota
)

// struct for a Tile
type Tile struct {
	X        int
	Y        int
	TileType int
}

// create a 2d array of Tile

func MakeDungeon() [][]Tile {
	dungeon := make([][]Tile, 10)

	// fill dungeon with TileType Open and set the X and Y from 0 to 10
	for y := 0; y < 10; y++ {
		dungeon[y] = make([]Tile, 10)

		for x := 0; x < 10; x++ {
			dungeon[y][x].TileType = Open
		}
	}
	return dungeon
}
