package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRobotAndCleaningSession(t *testing.T) {
	// Create a robot.
	robo1 := NewRobot("Johnny 5", 500)

	// Create a new rooms.
	area1 := NewArea("Work Room #1", 10000, 12000, 3)
	// area2 := NewArea("Work Room #2", 1000, 2000, 2)

	// Create a new cleaning session for our new
	// robot and room.
	robo1.NewCleaningSession(area1)

	require.Equal(t, "0.00", robo1.Session[0].Area[0].Completion(), "should have no squares completed")
	require.Equal(t, (area1.SizeX/robo1.Size)*(area1.SizeY/robo1.Size), len(robo1.Session[0].Area[0].Grid), "should have 480 squares total")

	// Clean the top most row on the grid.
	directionX := true
	xPos := 0
	passesIncreased := 0
	for i := 0; i < 120; i++ {
		if directionX {
			xPos += robo1.Size
		} else {
			xPos -= robo1.Size
		}
		if xPos >= area1.SizeX || xPos <= 0 {
			directionX = !directionX
		}
		inc := robo1.Session[0].Area[0].SetVisited(xPos+(robo1.Size/2), robo1.Size/2)
		if inc {
			passesIncreased++
		}
	}
	// area1.Print()
	// Dump(passesIncreased)
	require.Equal(t, "4.17", robo1.Session[0].Area[0].Completion(), "should be 4.17% completed (20/480 grid squares)")
}
