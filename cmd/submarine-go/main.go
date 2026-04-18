package main

import (
	"fmt"
	"log"

	"github.com/satokae/submarine-go/internal/agent"
	"github.com/satokae/submarine-go/internal/game"
	"github.com/satokae/submarine-go/internal/ui"
)

func main() {
	seed, err := ui.PromptSeed()
	if err != nil {
		log.Fatal(err)
	}

	isSecond, err := ui.PromptOrder()
	if err != nil {
		log.Fatal(err)
	}

	selectedAgent := agent.NewKashiwaAgent(seed)
	g := game.NewGame(selectedAgent, isSecond)

	if err := g.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
