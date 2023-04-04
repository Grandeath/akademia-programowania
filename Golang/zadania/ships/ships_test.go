package ships

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPoint(t *testing.T) {
	testCases := []struct {
		name        string
		point       Point
		addPoint    Point
		expectPoint Point
	}{
		{
			name:        "Zero case",
			point:       Point{0, 0},
			addPoint:    Point{0, 0},
			expectPoint: Point{0, 0},
		},
		{
			name:        "Add zero",
			point:       Point{1, 1},
			addPoint:    Point{0, 0},
			expectPoint: Point{1, 1},
		},
		{
			name:        "Add negative",
			point:       Point{1, 1},
			addPoint:    Point{-2, -1},
			expectPoint: Point{-1, 0},
		},
		{
			name:        "Add positive",
			point:       Point{1, 1},
			addPoint:    Point{2, 3},
			expectPoint: Point{3, 4},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newPoint := tc.point.Add(tc.addPoint)
			assert.Equal(t, tc.expectPoint.X, newPoint.X)
			assert.Equal(t, tc.expectPoint.Y, newPoint.Y)
		})
	}
}

func TestSize(t *testing.T) {
	testCases := []struct {
		name      string
		ships     Ship
		expectLen int
	}{
		{
			name:      "Null len",
			ships:     Ship{},
			expectLen: 0,
		},
		{
			name:      "Normal len",
			ships:     Ship{Point{1, 1}, Point{2, 2}},
			expectLen: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotLen := tc.ships.Size()
			assert.Equal(t, tc.expectLen, gotLen)
		})
	}
}

func TestMoveToPoint(t *testing.T) {
	testCases := []struct {
		name     string
		ship     Ship
		point    Point
		wantShip Ship
	}{
		{
			name:     "No movement",
			ship:     Ship{Point{1, 1}, Point{2, 2}},
			point:    Point{},
			wantShip: Ship{Point{0, 0}, Point{1, 1}},
		},
		{
			name:     "Normal movement",
			ship:     Ship{Point{1, 1}, Point{2, 2}},
			point:    Point{2, 2},
			wantShip: Ship{Point{2, 2}, Point{3, 3}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotShip := tc.ship.MoveTo(tc.point)
			assert.Equal(t, tc.wantShip, gotShip)
		})
	}
}
