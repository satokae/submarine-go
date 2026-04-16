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
