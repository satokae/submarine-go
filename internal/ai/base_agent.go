package ai

import (
	"math/rand/v2"

	"github.com/satokae/submarine-go/internal/constant"
	"github.com/satokae/submarine-go/internal/core"
)

func NewBaseAgent(name string, seed uint64, getInitialPositions func(rng *rand.Rand) []core.Position) *BaseAgent {
	r := rand.New(rand.NewPCG(seed, seed+1))
	positions := getInitialPositions(r)

	fleet := make(core.Fleet, len(positions))
	for i, pos := range positions {
		fleet[i] = core.Submarine{
			ID:       i,
			Position: pos,
			HP:       constant.InitialHP,
		}
	}
	return &BaseAgent{
		Name:           name,
		Rng:            r,
		FriendlyFleet:  fleet,
		SunkPositions:  []core.Position{},
		OffenseMap:     NewBeliefMap(constant.InitialSubmarines),
		DefenseMap:     NewBeliefMap(constant.InitialSubmarines),
		EnemyHPSum:     constant.InitialHP * constant.InitialSubmarines,
		WasHitLastTurn: false,
	}
}
