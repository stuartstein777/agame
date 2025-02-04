package entities

type Enemy struct {
	XY          Coord
	Health      int
	Attack      int
	Defense     int
	Speed       int
	SpriteIndex int
}
