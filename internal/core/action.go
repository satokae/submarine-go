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

func (o AttackOutcome) String() string {
	switch o {
	case Hit:
		return "命中"
	case HitAndSunk:
		return "命中・撃沈"
	case HighWaves:
		return "波高し"
	case Miss:
		return "ハズレ"
	default:
		return "不明"
	}
}

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
	MoveTarget   *Submarine
	AttackAction *AttackAction
}
