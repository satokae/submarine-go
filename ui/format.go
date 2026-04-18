package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
	"github.com/satokae/submarine-go/internal/ai"
	"github.com/satokae/submarine-go/internal/constant"
	"github.com/satokae/submarine-go/internal/core"
	"github.com/satokae/submarine-go/internal/util"
)

var (
	symbolEmpty        = color.HiBlackString(".")
	symbolSunkPos      = color.HiBlackString("x")
	symbolSubHP3       = color.GreenString("3")
	symbolSubHP2       = color.YellowString("2")
	symbolSubHP1       = color.RedString("1")
	symbolSubDestroyed = color.HiBlackString("X")
	symbolNeighbor     = color.MagentaString(".")
	symbolSeparator    = color.HiBlackString("│")
)

const (
	numberWidth = 5
)

func FormatBoard(name string, offenseMap ai.BeliefMap, defenseMap ai.BeliefMap, sunkPositions []core.Position, friendlyFleet core.Fleet) string {
	board := generateStatusBoard(sunkPositions, friendlyFleet)
	var output string

	output += "  1 2 3 4 5 " + symbolSeparator + " Enemy"
	output += strings.Repeat(" ", numberWidth*(constant.MapSize+2)-5)
	output += symbolSeparator + " "
	output += name
	output += "\n"

	for y := 0; y < constant.MapSize; y++ {
		output += string('A' + byte(y))
		for x := 0; x < constant.MapSize; x++ {
			output += " "
			output += board[0]
			board = board[1:]
		}
		output += " " + symbolSeparator + " "

		indexStart := y * constant.MapSize
		indexEnd := (y + 1) * constant.MapSize

		for _, val := range offenseMap.Grid()[indexStart:indexEnd] {
			colorFunc := getColorFunc(val)
			output += colorFunc(fmt.Sprintf("%*.*f", numberWidth, numberWidth-2, val))
			output += "  "
		}

		output += symbolSeparator + " "
		for _, val := range defenseMap.Grid()[indexStart:indexEnd] {
			colorFunc := getColorFunc(val)
			output += colorFunc(fmt.Sprintf("%*.*f", numberWidth, numberWidth-2, val))
			output += "  "
		}
		output += "\n"
	}
	return output
}

func generateStatusBoard(sunkPositions []core.Position, friendlyFleet core.Fleet) []string {
	board := make([]string, constant.GridSize)

	for i := range board {
		board[i] = symbolEmpty
	}

	for _, pos := range sunkPositions {
		board[pos] = symbolSunkPos
	}

	for _, sub := range friendlyFleet {
		var symbol string
		switch sub.HP {
		case 3:
			symbol = symbolSubHP3
		case 2:
			symbol = symbolSubHP2
		case 1:
			symbol = symbolSubHP1
		case 0:
			symbol = symbolSubDestroyed
		default:
			symbol = "?"
		}
		board[sub.Position] = symbol

		if sub.HP > 0 {
			neighbors := core.GetNeighbors(sub.Position)
			for _, n := range neighbors {
				if board[n] == symbolEmpty {
					board[n] = symbolNeighbor
				}
			}
		}
	}
	return board
}

func getRGBForValue(v float64) (int, int, int) {
	v = util.Clamp(v, 0, 1)
	gray := []int{128, 128, 128}
	yellow := []int{255, 255, 0}
	red := []int{255, 0, 0}

	var r, g, b int
	if v <= 0.5 {
		factor := v * 2
		r = int(math.Round(float64(gray[0]) + float64(yellow[0]-gray[0])*factor))
		g = int(math.Round(float64(gray[1]) + float64(yellow[1]-gray[1])*factor))
		b = int(math.Round(float64(gray[2]) + float64(yellow[2]-gray[2])*factor))
	} else {
		factor := (v - 0.5) * 2
		r = int(math.Round(float64(yellow[0]) + float64(red[0]-yellow[0])*factor))
		g = int(math.Round(float64(yellow[1]) + float64(red[1]-yellow[1])*factor))
		b = int(math.Round(float64(yellow[2]) + float64(red[2]-yellow[2])*factor))
	}
	return r, g, b
}

func getColorFunc(v float64) func(a ...interface{}) string {
	r, g, b := getRGBForValue(v)
	return color.RGB(r, g, b).SprintFunc()
}
