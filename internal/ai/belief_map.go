package ai

import (
	"github.com/satokae/submarine-go/internal/constant"
	"github.com/satokae/submarine-go/internal/core"
)

type BeliefMap struct {
	grid           []float64
	submarinesLeft int
	size           int
}

func NewBeliefMap(submarineCount int) *BeliefMap {
	grid := make([]float64, constant.GridSize)
	initialValue := float64(submarineCount) / float64(constant.GridSize)
	for i := range grid {
		grid[i] = initialValue
	}

	return &BeliefMap{
		grid:           grid,
		submarinesLeft: submarineCount,
		size:           constant.MapSize,
	}
}

func (m *BeliefMap) Grid() []float64 {
	return m.grid
}

func (m *BeliefMap) SubmarinesLeft() int {
	return m.submarinesLeft
}

func (m *BeliefMap) normalize(fixedPositions ...core.Position) {
	fixedIndexSet := make(map[int]bool)
	for _, pos := range fixedPositions {
		fixedIndexSet[int(pos)] = true
	}

	sumOfBeliefs := 0.0
	for i, val := range m.grid {
		if !fixedIndexSet[i] {
			sumOfBeliefs += val
		}
	}

	sumOfFixed := float64(len(fixedIndexSet))
	n := float64(m.submarinesLeft) - sumOfFixed

	if sumOfBeliefs == 0 {
		uniformValue := n / float64(len(m.grid)-len(fixedIndexSet))
		for i := range m.grid {
			if fixedIndexSet[i] {
				m.grid[i] = 1.0
			} else {
				m.grid[i] = uniformValue
			}
		}
		return
	}

	scale := n / sumOfBeliefs
	for i := range m.grid {
		if fixedIndexSet[i] {
			m.grid[i] = 1.0
		} else {
			m.grid[i] *= scale
		}
	}
}

func (m *BeliefMap) UpdateOnHit(pos core.Position) {
	m.grid[pos] = 1.0
	m.normalize(pos)
}

func (m *BeliefMap) UpdateOnSunk(pos core.Position) {
	m.grid[pos] = 0.0
	m.submarinesLeft--
	m.normalize()
}

func (m *BeliefMap) UpdateOnMove(direction core.Direction, distance int, sunkPositions []core.Position) {
	n := float64(m.submarinesLeft)
	if n <= 0 {
		return
	}

	newGrid := make([]float64, len(m.grid))
	dx, dy := direction.ToVector()

	isSunk := func(p core.Position) bool {
		for _, sp := range sunkPositions {
			if sp == p {
				return true
			}
		}
		return false
	}

	for i, val := range m.grid {
		if val == 0 {
			continue
		}

		movingPart := val * (1.0 / n)
		stationaryPart := val * ((n - 1.0) / n)

		newGrid[i] += stationaryPart

		moveTo, err := core.Position(i).Move(dx*distance, dy*distance)
		isPlausible := err == nil && !isSunk(moveTo)
		if distance == 2 {
			intermediatePos, err := core.Position(i).Move(dx, dy)
			if err != nil || isSunk(intermediatePos) {
				isPlausible = false
			}
		}

		if isPlausible {
			newGrid[moveTo] += movingPart
		} else {
			newGrid[i] += movingPart
		}
	}

	m.grid = newGrid
}
