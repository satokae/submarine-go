package core_test

import (
	"testing"

	"github.com/satokae/submarine-go/internal/core"
)

func TestGetValidMoveDestination(t *testing.T) {
	fleet := core.Fleet{
		{ID: 0, Position: 0, HP: 3}, // Top-left (0,0)
		{ID: 1, Position: 2, HP: 3}, // (2,0)
	}
	sunk := []core.Position{5} // (0,1)

	tests := []struct {
		name      string
		move      core.MoveAction
		target    core.Submarine
		wantPos   core.Position
		wantValid bool
	}{
		{
			name:      "Valid Move: East 1",
			move:      core.MoveAction{Direction: core.East, Distance: 1},
			target:    fleet[0],
			wantPos:   1,
			wantValid: true,
		},
		{
			name:      "Invalid Move: Destination Occupied by another sub",
			move:      core.MoveAction{Direction: core.East, Distance: 2},
			target:    fleet[0],
			wantPos:   0, // Returns original pos on failure
			wantValid: false,
		},
		{
			name:      "Invalid Move: Out of Bounds",
			move:      core.MoveAction{Direction: core.North, Distance: 1},
			target:    fleet[0],
			wantPos:   0,
			wantValid: false,
		},
		{
			name:      "Invalid Move: Intermediate Occupied (Sunk Submarine)",
			move:      core.MoveAction{Direction: core.South, Distance: 2},
			target:    fleet[0],
			wantPos:   0,
			wantValid: false,
		},
		{
			name:      "Valid Move: South 1",
			move:      core.MoveAction{Direction: core.South, Distance: 1},
			target:    fleet[1],
			wantPos:   7,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPos, gotValid := core.GetValidMoveDestination(tt.move, tt.target, fleet, sunk)
			if gotValid != tt.wantValid {
				t.Errorf("%s: gotValid = %v, want %v", tt.name, gotValid, tt.wantValid)
			}
			if gotPos != tt.wantPos {
				t.Errorf("%s: gotPos = %v, want %v", tt.name, gotPos, tt.wantPos)
			}
		})
	}
}

func TestIsAttackPossible(t *testing.T) {
	fleet := core.Fleet{
		{ID: 0, Position: 12, HP: 3}, // Center (2,2)
		{ID: 1, Position: 0, HP: 0},  // Sunk (0,0)
	}

	tests := []struct {
		name   string
		attack core.AttackAction
		want   bool
	}{
		{"Possible: Neighbor is active sub", core.AttackAction{Position: 13}, true},
		{"Impossible: Target is active sub itself", core.AttackAction{Position: 12}, false},
		{"Impossible: Neighbor is only a sunk sub", core.AttackAction{Position: 1}, false},
		{"Impossible: No active sub nearby", core.AttackAction{Position: 24}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.IsAttackPossible(tt.attack, fleet); got != tt.want {
				t.Errorf("%s: IsAttackPossible() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
