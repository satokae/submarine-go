package agent

import (
	"math"
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

func (k *KashiwaAgent) ChooseAction() *core.Action {
	possibleAttacks := k.GenerateAllPossibleAttacks()
	possibleMoves := k.GenerateAllPossibleMoves()

	bestAttack := k.findBestAttack(possibleAttacks, len(possibleMoves) > 0)
	if bestAttack != nil {
		return &core.Action{
			Type:         core.ActionTypeAttack,
			AttackAction: bestAttack,
		}
	}

	bestMove, moveTarget := k.findBestMove(possibleMoves)
	if bestMove != nil && moveTarget != nil {
		return &core.Action{
			Type:       core.ActionTypeMove,
			MoveTarget: moveTarget,
			MoveAction: bestMove,
		}
	}

	return nil
}

func (k *KashiwaAgent) calculateAttackScore(attack core.AttackAction, hasMoves bool) float64 {
	offenseValue := k.OffenseMap.Grid()[attack.Position]

	if (offenseValue < 0.05 && hasMoves) || (k.WasHitLastTurn && offenseValue < 0.85 && hasMoves) {
		return math.Inf(-1)
	}

	if offenseValue < 0.75 {
		currentEntropy := k.OffenseMap.CalculateEntropy()
		tempMap := k.OffenseMap.Copy()
		tempMap.UpdateOnNear(attack.Position)
		entropyOnNear := tempMap.CalculateEntropy()
		return currentEntropy - entropyOnNear
	} else {
		return offenseValue * 1.5
	}
}

func (k *KashiwaAgent) calculateMoveScore(action core.Action) float64 {
	moveFrom := action.MoveTarget.Position
	if _, ok := core.GetValidMoveDestination(*action.MoveAction, *action.MoveTarget, k.FriendlyFleet, k.SunkPositions); !ok {
		return math.Inf(-1)
	}
	return k.DefenseMap.Grid()[moveFrom] + k.Rng.Float64()*8
}

func (k *KashiwaAgent) findBestAttack(attacks []core.AttackAction, hasMoves bool) *core.AttackAction {
	if len(attacks) == 0 {
		return nil
	}

	type scoredAttack struct {
		attack core.AttackAction
		score  float64
	}

	scoredAttacks := []scoredAttack{}
	maxScore := math.Inf(-1)

	for _, attack := range attacks {
		score := k.calculateAttackScore(attack, hasMoves)
		if score > maxScore {
			maxScore = score
		}
		scoredAttacks = append(scoredAttacks, scoredAttack{attack, score})
	}

	if math.IsInf(maxScore, -1) {
		return nil
	}

	bestAttacks := []core.AttackAction{}
	for _, sa := range scoredAttacks {
		if sa.score == maxScore {
			bestAttacks = append(bestAttacks, sa.attack)
		}
	}

	res := bestAttacks[k.Rng.IntN(len(bestAttacks))]
	return &res
}

func (k *KashiwaAgent) findBestMove(moves []core.Action) (*core.MoveAction, *core.Submarine) {
	if len(moves) == 0 {
		return nil, nil
	}

	type scoredMove struct {
		move  core.Action
		score float64
	}

	scoredMoves := []scoredMove{}
	maxScore := math.Inf(-1)

	for _, move := range moves {
		score := k.calculateMoveScore(move)
		if score > maxScore {
			maxScore = score
		}
		scoredMoves = append(scoredMoves, scoredMove{move, score})
	}

	if math.IsInf(maxScore, -1) {
		return nil, nil
	}

	bestMoves := []core.Action{}
	for _, sm := range scoredMoves {
		if sm.score == maxScore {
			bestMoves = append(bestMoves, sm.move)
		}
	}

	res := bestMoves[k.Rng.IntN(len(bestMoves))]
	return res.MoveAction, res.MoveTarget
}
