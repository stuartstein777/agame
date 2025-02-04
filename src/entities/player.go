package entities

type Player struct {
	PlayLocation Coord
	Health       int
	Attack       int
	Defense      int
	SpriteIndex  int
	Items        []Item
}
