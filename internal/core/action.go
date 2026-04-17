package core

type AttackOutcome int
type ActionType int

const (
	Hit AttackOutcome = iota
	HitAndSunk
	HighWaves
	Miss
)

const (
	ActionTypeMove ActionType = iota
	ActionTypeAttack
)

type MoveAction struct {
	Direction Direction
	Distance  int
}

type AttackAction struct {
	Position Position
}

type Action struct {
	Type         ActionType
	MoveAction   *MoveAction
	AttackAction *AttackAction
}
