package agent

import (
	"math/rand/v2"

	"github.com/satokae/submarine-go/internal/ai"
	"github.com/satokae/submarine-go/internal/constant"
	"github.com/satokae/submarine-go/internal/core"
)

var _ ai.Agent = (*KashiwaAgent)(nil)

const minCoverageThreshold = 21

type KashiwaAgent struct {
	*ai.BaseAgent
}

func NewKashiwaAgent(seed uint64) *KashiwaAgent {
	getInitialPositions := func(rng *rand.Rand) []core.Position {
		for {
			indexSet := make(map[int]bool)
			for len(indexSet) < constant.InitialSubmarines {
				indexSet[rng.IntN(constant.GridSize)] = true
			}

			positions := []core.Position{}
			for idx := range indexSet {
				positions = append(positions, core.Position(idx))
			}

			positionSet := make(map[int]bool)
			for _, pos := range positions {
				positionSet[int(pos)] = true
				neighbors := core.GetNeighbors(pos)
				for _, n := range neighbors {
					positionSet[int(n)] = true
				}
			}

			if len(positionSet) > minCoverageThreshold {
				return positions
			}
		}
	}

	return &KashiwaAgent{
		BaseAgent: ai.NewBaseAgent("Kashiwa", seed, getInitialPositions),
	}
}
