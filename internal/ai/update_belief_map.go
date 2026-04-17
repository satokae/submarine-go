package ai

import "github.com/satokae/submarine-go/internal/core"

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

func (m *BeliefMap) UpdateOnMiss(pos core.Position) {
	neighbors := core.GetNeighbors(pos)
	m.grid[pos] = 0.0

	for _, n := range neighbors {
		m.grid[n] = 0.0
	}
	m.normalize()
}
