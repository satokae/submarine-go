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

func (b *BaseAgent) GetEnemyHPSum() int {
	return b.EnemyHPSum
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

func (b *BaseAgent) GenerateAllPossibleMoves() []core.Action {
	moves := []core.Action{}
	distances := []int{1, 2}
	directions := []core.Direction{core.East, core.North, core.South, core.West}

	for i := range b.AvailableFleet() {
		sub := &b.FriendlyFleet[i]
		for _, dir := range directions {
			for _, dist := range distances {
				move := core.Action{
					Type:       core.ActionTypeMove,
					MoveTarget: sub,
					MoveAction: &core.MoveAction{
						Direction: dir,
						Distance:  dist,
					},
				}

				if _, ok := core.GetValidMoveDestination(*move.MoveAction, *sub, b.FriendlyFleet, b.SunkPositions); ok {
					moves = append(moves, move)
				}
			}
		}
	}
	return moves
}

func (b *BaseAgent) GenerateAllPossibleAttacks() []core.AttackAction {
	attacks := []core.AttackAction{}

	for i := 0; i < constant.GridSize; i++ {
		attack := core.AttackAction{Position: core.Position(i)}
		if !core.IsAttackPossible(attack, b.FriendlyFleet) {
			continue
		}

		attacks = append(attacks, attack)
	}

	return attacks
}

func (b *BaseAgent) ApplyOutcomeToMap(m *BeliefMap, pos core.Position, outcome core.AttackOutcome) {
	switch outcome {
	case core.Hit:
		m.UpdateOnHit(pos)
	case core.HitAndSunk:
		m.UpdateOnSunk(pos)
	case core.HighWaves:
		m.UpdateOnNear(pos)
	case core.Miss:
		m.UpdateOnMiss(pos)
	}
}

func (b *BaseAgent) ApplyActionToMap(m *BeliefMap, action core.Action) {
	if action.Type == core.ActionTypeMove {
		m.UpdateOnMove(action.MoveAction.Direction, action.MoveAction.Distance, b.SunkPositions)
	} else {
		m.UpdateOnNear(action.AttackAction.Position)
	}
}

func (b *BaseAgent) OnAttackResult(pos core.Position, outcome core.AttackOutcome) {
	b.ApplyOutcomeToMap(b.OffenseMap, pos, outcome)

	if outcome == core.Hit {
		b.EnemyHPSum--
	}

	if outcome == core.HitAndSunk {
		b.EnemyHPSum--
		b.SunkPositions = append(b.SunkPositions, pos)
	}
}

func (b *BaseAgent) OnDefenseResult(pos core.Position, outcome core.AttackOutcome) {
	b.ApplyOutcomeToMap(b.DefenseMap, pos, outcome)
	b.WasHitLastTurn = (outcome == core.Hit)

	if outcome == core.Hit || outcome == core.HitAndSunk {
		for i := range b.FriendlyFleet {
			if b.FriendlyFleet[i].Position == pos {
				b.FriendlyFleet[i].HP--
				if outcome == core.HitAndSunk {
					b.SunkPositions = append(b.SunkPositions, pos)
				}
				break
			}
		}
	}
}

func (b *BaseAgent) OnOwnAction(action core.Action) {
	b.ApplyActionToMap(b.DefenseMap, action)

	if action.Type == core.ActionTypeMove {
		target := *action.MoveTarget
		dest, ok := core.GetValidMoveDestination(*action.MoveAction, target, b.FriendlyFleet, b.SunkPositions)
		if ok {
			for i := range b.FriendlyFleet {
				if b.FriendlyFleet[i].ID == target.ID {
					b.FriendlyFleet[i].Position = dest
					break
				}
			}
		}
	}
}

func (b *BaseAgent) OnEnemyAction(action core.Action) {
	b.WasHitLastTurn = false
	b.ApplyActionToMap(b.OffenseMap, action)
}

func (b *BaseAgent) AgentName() string {
	return b.Name
}

func (b *BaseAgent) Fleet() core.Fleet {
	fleet := make(core.Fleet, len(b.FriendlyFleet))
	copy(fleet, b.FriendlyFleet)
	return fleet
}

func (b *BaseAgent) GetOffenseMap() *BeliefMap {
	return b.OffenseMap
}

func (b *BaseAgent) GetDefenseMap() *BeliefMap {
	return b.DefenseMap
}

func (b *BaseAgent) GetSunkPositions() []core.Position {
	return b.SunkPositions
}
