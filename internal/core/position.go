package core

import (
	"errors"

	"github.com/satokae/submarine-go/internal/constant"
)

type Position int

var ErrOutOfBounds = errors.New("position out of bounds")

func (pos Position) Move(dx int, dy int) (Position, error) {
	x := int(pos) % constant.MapSize
	y := int(pos) / constant.MapSize

	x += dx
	y += dy
	if x < 0 || x >= constant.MapSize || y < 0 || y >= constant.MapSize {
		return -1, (ErrOutOfBounds)
	}
	return Position(y*constant.MapSize + x), nil
}

func GetNeighbors(pos Position) []Position {
	neighbors := []Position{}

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			neighborPos, err := pos.Move(dx, dy)
			if err == nil {
				neighbors = append(neighbors, neighborPos)
			}
		}
	}
	return neighbors
}
