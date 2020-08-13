package entity

import (
	"fmt"
	"time"
)

// Square represents a grid square to be cleaned.
type Square struct {
	X               int        `json:"x,omitempty"`
	Y               int        `json:"y,omitempty"`
	Size            int        `json:"size,omitempty"`
	Passes          int        `json:"passes,omitempty"`
	HasRobotPresent bool       `json:"has_robot_present,omitempty"` // Is the robot currently on the square?
	CleanedAt       *time.Time `json:"cleaned_at,omitempty"`

	// To ensure we can retrieve all grid squares in the order they were created.
	Order int `json:"order,omitempty"`

	Common
}

// IsInSquare checks if the given x,y coordinates fall within
// the square.
func (s *Square) IsInSquare(x, y int) bool {
	isInX := s.X <= x && x < (s.X+s.Size)
	isInY := s.Y <= y && y < (s.Y+s.Size)
	return isInX && isInY
}

// NewGrid creates a new grid given an area defined by x and y.
func NewGrid(xmm, ymm, size int) []*Square {
	var s []*Square
	var count int
	for y := 0; y < ymm; y += size {
		for x := 0; x < xmm; x += size {
			count++
			s = append(s, &Square{
				X:     x,
				Y:     y,
				Size:  size,
				Order: count,
				Common: Common{
					UID:       fmt.Sprintf("_:%s%d", SquareUID, count),
					DType:     []string{"Square"},
					CreatedAt: now(),
				},
			})
		}
	}
	return s
}
