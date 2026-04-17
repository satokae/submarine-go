package core_test

import (
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/satokae/submarine-go/internal/core"
)

// TestPosition_Move validates both successful coordinate translation
// and out-of-bounds error handling.
func TestPosition_Move(t *testing.T) {
	tests := []struct {
		name    string
		pos     core.Position
		dx, dy  int
		wantPos core.Position
		wantErr error
	}{
		{"Valid: Move East", 0, 1, 0, 1, nil},
		{"Valid: Move South", 0, 0, 1, 5, nil},
		{"Valid: Move Diagonal", 0, 1, 1, 6, nil},
		{"Valid: Far move", 0, 4, 4, 24, nil},
		{"Error: North OOB", 2, 0, -1, -1, core.ErrOutOfBounds},
		{"Error: West OOB", 5, -1, 0, -1, core.ErrOutOfBounds},
		{"Error: East OOB", 4, 1, 0, -1, core.ErrOutOfBounds},
		{"Error: South OOB", 20, 0, 5, -1, core.ErrOutOfBounds},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPos, err := tt.pos.Move(tt.dx, tt.dy)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPos != tt.wantPos {
				t.Errorf("Move() gotPos = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}

// TestGetNeighbors ensures that we get the correct Moore neighborhood
// while correctly filtering out tiles off the grid.
func TestGetNeighbors(t *testing.T) {
	tests := []struct {
		name string
		pos  core.Position
		want []core.Position
	}{
		{
			name: "Top-Left Corner",
			pos:  0,
			want: []core.Position{1, 5, 6},
		},
		{
			name: "Center Tile",
			pos:  12,
			want: []core.Position{6, 7, 8, 11, 13, 16, 17, 18},
		},
		{
			name: "Bottom-Right Corner",
			pos:  24,
			want: []core.Position{18, 19, 23},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.GetNeighbors(tt.pos)

			sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
			sort.Slice(tt.want, func(i, j int) bool { return tt.want[i] < tt.want[j] })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNeighbors(%v) = %v, want %v", tt.pos, got, tt.want)
			}
		})
	}
}
