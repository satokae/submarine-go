package ai

import (
	"math"

	"github.com/satokae/submarine-go/internal/core"
)

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

func (m *BeliefMap) UpdateOnNear(pos core.Position) {
	type region int

	const (
		far region = iota
		near
		center
	)

	n := m.submarinesLeft
	if n <= 0 {
		return
	}

	neighbors := core.GetNeighbors(pos)
	likelihood := make([]float64, len(m.grid))
	regionKind := make([]region, len(m.grid))
	sumR := 0.0

	neighborSet := make(map[core.Position]bool)
	for _, n := range neighbors {
		neighborSet[n] = true
	}

	for i := range regionKind {
		p := core.Position(i)
		if p == pos {
			regionKind[i] = center
		} else if neighborSet[p] {
			regionKind[i] = near
		} else {
			regionKind[i] = far
			sumR += m.grid[i]
		}
	}

	if n == 1 {
		for i := range m.grid {
			if regionKind[i] == near {
				likelihood[i] = 1.0
			} else {
				likelihood[i] = 0.0

			}
		}
	} else {
		// 1 < n
		eC := m.grid[pos]
		decrementedN := float64(n - 1)

		for i := range m.grid {
			eI := m.grid[i]
			denom := float64(n) - eI

			if denom <= 0 || decrementedN <= 0 {
				likelihood[i] = 0.0
				continue
			}

			switch regionKind[i] {
			case center:
				likelihood[i] = 0.0

			case near:
				likelihood[i] = math.Pow(1.0-eC/denom, decrementedN)

			case far:
				eCPrime := (eC / denom) * decrementedN
				sumRPrime := ((sumR - eI) / denom) * decrementedN

				pNoCGivenIDenom := decrementedN
				pNoCGivenI := 0.0
				if pNoCGivenIDenom > 0 {
					pNoCGivenI = math.Pow(1.0-eCPrime/pNoCGivenIDenom, decrementedN)
				}

				pAllInROrCGivenIBase := sumRPrime + eCPrime
				pAllInROrCGivenIDenom := decrementedN
				pAllInROrCGivenI := 0.0
				if pAllInROrCGivenIDenom > 0 {
					pAllInROrCGivenI = math.Pow(pAllInROrCGivenIBase/pAllInROrCGivenIDenom, decrementedN)
				}

				pAtLeastOneInSGivenI := 1.0 - pAllInROrCGivenI
				likelihood[i] = pNoCGivenI * pAtLeastOneInSGivenI
			}
		}
	}

	for i := range m.grid {
		m.grid[i] *= likelihood[i]
	}
	m.normalize()
}

func (m *BeliefMap) UpdateOnMiss(pos core.Position) {
	neighbors := core.GetNeighbors(pos)
	m.grid[pos] = 0.0

	for _, n := range neighbors {
		m.grid[n] = 0.0
	}
	m.normalize()
}
