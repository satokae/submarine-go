package game

import (
	"fmt"

	"github.com/satokae/submarine-go/internal/ai"
	"github.com/satokae/submarine-go/internal/constant"
	"github.com/satokae/submarine-go/internal/core"
	"github.com/satokae/submarine-go/internal/ui"
)

type Game struct {
	AIAgent        ai.Agent
	IsPlayerTurn   bool
	PlayerHPSum    int
	PlayerSubsLeft int
	TurnCount      int
}

func NewGame(agent ai.Agent, isPlayerTurn bool) *Game {
	return &Game{
		AIAgent:        agent,
		IsPlayerTurn:   isPlayerTurn,
		PlayerHPSum:    constant.InitialSubmarines * constant.InitialHP,
		PlayerSubsLeft: constant.InitialSubmarines,
		TurnCount:      1,
	}
}

func (g *Game) Start() error {
	for {
		fmt.Printf("[ Turn %d ]\n", g.TurnCount)
		fmt.Print(ui.FormatBoard(g.AIAgent))

		if g.AIAgent.GetEnemyHPSum() <= 0 {
			fmt.Print("相手潜水艦が全て撃沈されました\n")
			break
		}
		if g.AIAgent.OwnHPSum() <= 0 {
			fmt.Printf("%s の潜水艦が全て撃沈されました\n", g.AIAgent.AgentName())
		}

		if g.IsPlayerTurn {
			err := g.playerTurn()
			if err != nil {
				return err
			}
		} else {
			err := g.aiTurn()
			if err != nil {
				return err
			}
		}

		g.IsPlayerTurn = !g.IsPlayerTurn
		g.TurnCount++
		fmt.Println()
	}
	return nil
}

func (g *Game) playerTurn() error {
	actionType, err := ui.PromptActionType()
	if err != nil {
		return err
	}

	switch actionType {
	case core.ActionTypeAttack:
		pos, err := ui.PromptAttackPosition()
		if err != nil {
			return err
		}

		outcome, _ := core.ResolveAttack(core.AttackAction{Position: pos}, g.AIAgent.Fleet())
		fmt.Printf("[%s] 結果: %v\n", g.AIAgent.AgentName(), outcome)
		action := core.Action{
			Type:         core.ActionTypeAttack,
			AttackAction: &core.AttackAction{Position: pos},
		}

		g.AIAgent.OnEnemyAction(action)
		g.AIAgent.OnDefenseResult(pos, outcome)

	case core.ActionTypeMove:
		dir, err := ui.PromptMoveDirection()
		if err != nil {
			return err
		}
		dist, err := ui.PromptMoveDistance()
		if err != nil {
			return err
		}
		action := core.Action{
			Type: core.ActionTypeMove,
			MoveAction: &core.MoveAction{
				Direction: dir,
				Distance:  dist,
			},
		}

		g.AIAgent.OnEnemyAction(action)
	}

	return nil
}

func (g *Game) aiTurn() error {
	action := g.AIAgent.ChooseAction()
	if action == nil {
		fmt.Printf("[%s] 行動を選択できませんでした\n", g.AIAgent.AgentName())
		return nil
	}

	switch action.Type {
	case core.ActionTypeAttack:
		pos := action.AttackAction.Position

		fmt.Printf("[%s] Attack %v\n", g.AIAgent.AgentName(), pos)

		g.AIAgent.OnOwnAction(*action)

		outcome, err := ui.PromptAttackOutcome()
		if err != nil {
			return err
		}
		g.AIAgent.OnAttackResult(pos, outcome)

		if outcome == core.Hit || outcome == core.HitAndSunk {
			g.PlayerHPSum--
			if outcome == core.HitAndSunk {
				g.PlayerSubsLeft--
			}
		}

	case core.ActionTypeMove:
		move := action.MoveAction
		fmt.Printf("[%s] Move %v %d\n", g.AIAgent.AgentName(), move.Direction, move.Distance)
		g.AIAgent.OnOwnAction(*action)
	}
	return nil
}
