package entities

const (
	Key = iota
)

type Item struct {
	Type           int
	SpritesheetIdx int
}

// struct for a Tile
type Tile struct {
	X        int
	Y        int
	TileType int
	Items    []Item
}
