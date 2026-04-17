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

func (b *BaseAgent) OwnHPSum() int {
	sum := 0
	for _, sub := range b.FriendlyFleet {
		sum += sub.HP
	}
	return sum
}

func (b *BaseAgent) AvailableFleet() core.Fleet {
	available := []core.Submarine{}
	for _, sub := range b.FriendlyFleet {
		if sub.HP > 0 {
			available = append(available, sub)
		}
	}
	return available
}

func (b *BaseAgent) GenerateAllPossibleMoves() []core.MoveAction {
	moves := []core.MoveAction{}
	distances := []int{1, 2}
	directions := []core.Direction{core.East, core.North, core.South, core.West}

	for i := range b.AvailableFleet() {
		sub := &b.FriendlyFleet[i]
		for _, dir := range directions {
			for _, dist := range distances {
				move := core.MoveAction{
					Direction: dir,
					Distance:  dist,
				}

				if _, ok := core.GetValidMoveDestination(move, *sub, b.FriendlyFleet, b.SunkPositions); ok {
					moves = append(moves, move)
				}
			}
		}
	}
	return moves
}
