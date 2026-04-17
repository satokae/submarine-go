package ai

import (
	"math"

	"github.com/satokae/submarine-go/internal/constant"
	"github.com/satokae/submarine-go/internal/core"
	"github.com/satokae/submarine-go/internal/util"
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

func (m *BeliefMap) CalculateEntropy() float64 {
	entropy := 0.0
	for _, p := range m.grid {
		val := util.Clamp(p, 0, 1)

		if 0 < val && val < 1 {
			e := -val*math.Log2(val) - (1.0-val)*math.Log2(1.0-val)
			entropy += e
		}

	}
	return entropy
}

func (m *BeliefMap) Copy() *BeliefMap {
	newGrid := make([]float64, len(m.grid))
	copy(newGrid, m.grid)
	return &BeliefMap{
		grid:           newGrid,
		submarinesLeft: m.submarinesLeft,
		size:           m.size,
	}
}
