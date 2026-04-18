package game

import (
	"github.com/satokae/submarine-go/internal/ai"
	"github.com/satokae/submarine-go/internal/constant"
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
