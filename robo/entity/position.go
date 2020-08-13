package entity

import (
	"time"
)

// Position represents a robot's position within
// an area at a certain time.
type Position struct {
	X        int        `json:"x,omitempty"`
	Y        int        `json:"y,omitempty"`
	PassedAt *time.Time `json:"passed_at,omitempty"`
	Common
}

// NewPosition creates a new position.
func NewPosition(x int, y int, passedAt time.Time) *Position {
	return &Position{
		X:        x,
		Y:        y,
		PassedAt: &passedAt,
		Common: Common{
			UID:       "_:" + PositionUID,
			DType:     []string{"Position"},
			CreatedAt: now(),
		},
	}
}
