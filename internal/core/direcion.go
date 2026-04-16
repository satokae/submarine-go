package core

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

func (d Direction) ToVector() (dx, dy int) {
	switch d {
	case North:
		return 0, -1
	case South:
		return 0, 1
	case East:
		return 1, 0
	case West:
		return -1, 0
	default:
		return 0, 0
	}
}
