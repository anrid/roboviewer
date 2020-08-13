package entity

import (
	"fmt"
	"log"
	"math"
	"time"
)

// CleaningArea is a grid based on an area.
// It holds historical cleaning information.
type CleaningArea struct {
	Name  string `json:"name,omitempty"`   // Each cleaning area should definitely have a name to make reports nicer.
	SizeX int    `json:"size_x,omitempty"` // X side size in millimeters.
	SizeY int    `json:"size_y,omitempty"` // Y side size in millimeters.

	// The size of a grid square. Typically the same size os the diameter of the assigned cleaning robot.
	Grid []*Square `json:"grid,omitempty"`

	// Number of grid square passes needed before the square can be considered clean.
	PassesNeeded int `json:"passes_needed,omitempty"`

	Common
}

// NewCleaningArea creates a new cleaning area with a grid
// based on the given area and robot that's going to do the
// cleaning.
func NewCleaningArea(a *Area, r *Robot) *CleaningArea {
	if a.SizeX == 0 || a.SizeY == 0 || r.Size == 0 {
		log.Panicf(
			"failed to create a new cleaning area with invalid size values: area.x = %d area.y = %d robot.size = %d",
			a.SizeX, a.SizeX, r.Size,
		)
	}
	return &CleaningArea{
		Name:         a.Name,
		SizeX:        a.SizeX,
		SizeY:        a.SizeY,
		PassesNeeded: a.PassesNeeded,
		Grid:         NewGrid(a.SizeX, a.SizeY, r.Size),
		Common: Common{
			UID:       "_:" + CleaningAreaUID,
			DType:     []string{"CleaningArea"},
			CreatedAt: now(),
		},
	}
}

// Completion returns the completion percentage for an area
// as as 2-decimal string.
func (a *CleaningArea) Completion() string {
	var cleaned int
	for _, s := range a.Grid {
		if s.CleanedAt != nil {
			cleaned++
		}
	}
	if len(a.Grid) == 0 {
		return "0.00"
	}
	pct := float64(cleaned) / float64(len(a.Grid))
	return fmt.Sprintf("%0.2f", math.Round(pct*10000)/100)
}

// SetVisited marks a grid square as having been visited by
// the robot.
func (a *CleaningArea) SetVisited(x, y int) bool {
	now := time.Now()
	var registeredPass bool
	for _, s := range a.Grid {
		if s.IsInSquare(x, y) {
			if !s.HasRobotPresent {
				// Increase the number of passes if the robot
				// just entered this square.
				s.Passes++
				if s.Passes == a.PassesNeeded {
					s.CleanedAt = &now
				}
				s.HasRobotPresent = true
				registeredPass = true
			}
		} else {
			// Unlock this square since the robot is not
			// present.
			s.HasRobotPresent = false
		}
	}
	return registeredPass
}

// Print prints an ASCII representation of the grid and
// it's progress.
func (a *CleaningArea) Print() {
	row := 0
	for i, s := range a.Grid {
		if s.CleanedAt != nil {
			print("*")
		} else if s.Passes > 0 {
			print(s.Passes)
		} else {
			print("_")
		}
		if i == len(a.Grid)-1 {
			row++
			print(" ", row)
			println("")
		}
		if i <= len(a.Grid)-2 {
			if a.Grid[i+1].Y > s.Y {
				row++
				print(" ", row)
				println("")
			}
		}
	}
}
