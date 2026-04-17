package ai

import (
	"math/rand/v2"

	"github.com/satokae/submarine-go/internal/core"
)

type Agent interface {
	ChooseAction() *core.Action
	OnAttackResult(pos core.Position, outcome core.AttackOutcome)
	OnDefenseResult(pos core.Position, outcome core.AttackOutcome)
	OnOwnAction(action core.Action)
	OnEnemyAction(action core.Action)
	AvailableFleet() core.Fleet
	OwnHPSum() int
}

type BaseAgent struct {
	Name           string
	Rng            *rand.Rand
	FriendlyFleet  core.Fleet
	SunkPositions  []core.Position
	OffenseMap     *BeliefMap
	DefenseMap     *BeliefMap
	EnemyHPSum     int
	WasHitLastTurn bool
}
