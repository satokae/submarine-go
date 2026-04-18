package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/satokae/submarine-go/internal/core"
)

func generateSelectTemplate(name string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "  > " + name + ": {{ . }}",
	}
}

func generateInputTemplate(name string) *promptui.PromptTemplates {
	return &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   `{{ "  o" | green }} {{ . }}: `,
		Invalid: `{{ "  x" | red }} {{ . }}: `,
		Success: "  > " + name + ": ",
	}
}

type selectOption[T any] struct {
	Label string
	Value T
}

func runSelectPrompt[T any](label, templateName string, options []selectOption[T]) (T, error) {
	items := make([]string, len(options))
	for i, opt := range options {
		items[i] = opt.Label
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: generateSelectTemplate(templateName),
	}

	i, _, err := prompt.Run()
	if err != nil {
		var zero T
		return zero, err
	}

	return options[i].Value, nil
}

func PromptActionType() (core.ActionType, error) {
	return runSelectPrompt(
		"行動の種類を選択",
		"行動",
		[]selectOption[core.ActionType]{
			{Label: "攻撃", Value: core.ActionTypeAttack},
			{Label: "移動", Value: core.ActionTypeMove},
		},
	)
}

func PromptAttackOutcome() (core.AttackOutcome, error) {
	return runSelectPrompt(
		"攻撃結果を入力",
		"結果",
		[]selectOption[core.AttackOutcome]{
			{Label: "命中", Value: core.Hit},
			{Label: "命中・撃沈", Value: core.HitAndSunk},
			{Label: "波高し", Value: core.HighWaves},
			{Label: "ハズレ", Value: core.Miss},
		},
	)
}

func PromptMoveDirection() (core.Direction, error) {
	return runSelectPrompt(
		"移動方向を選択",
		"方向",
		[]selectOption[core.Direction]{
			{Label: "North", Value: core.North},
			{Label: "South", Value: core.South},
			{Label: "West", Value: core.West},
			{Label: "East", Value: core.East},
		},
	)
}

func PromptOrder() (bool, error) {
	return runSelectPrompt(
		"行動順を選択",
		"行動順",
		[]selectOption[bool]{
			{Label: "先攻", Value: false},
			{Label: "後攻", Value: true},
		},
	)
}

func PromptMoveDistance() (int, error) {
	return runSelectPrompt(
		"移動距離を選択",
		"距離",
		[]selectOption[int]{
			{Label: "1", Value: 1},
			{Label: "2", Value: 2},
		},
	)
}

func PromptSeed() (uint64, error) {
	validate := func(input string) error {
		var seed uint64
		_, err := fmt.Sscanf(input, "%d", &seed)
		if err != nil {
			return fmt.Errorf("数値を入力してください")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:     "シード",
		Validate:  validate,
		Templates: generateInputTemplate("シード"),
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	var seed uint64
	fmt.Sscanf(result, "%d", &seed)
	return seed, nil
}

func PromptAttackPosition() (core.Position, error) {
	validate := func(input string) error {
		if len(input) != 2 {
			return fmt.Errorf("A1からE5の形式で入力してください")
		}
		y := strings.ToUpper(input)[0]
		x := input[1]
		if y < 'A' || y > 'E' || x < '1' || x > '5' {
			return fmt.Errorf("A1からE5の形式で入力してください")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:     "攻撃座標",
		Validate:  validate,
		Templates: generateInputTemplate("座標"),
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	y := int(strings.ToUpper(result)[0] - 'A')
	x := int(result[1] - '1')
	return core.Position(y*5 + x), nil
}
