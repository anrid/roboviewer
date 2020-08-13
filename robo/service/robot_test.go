package service

import (
	"context"
	"testing"
	"time"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// RobotTestSuite defines the test suite.
type RobotTestSuite struct {
	suite.Suite
	ctx context.Context
	th  *testHelper
}

func (s *RobotTestSuite) TestListRobots() {
	robots, err := s.th.Service.Robot.List(s.ctx, "", "")
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(robots))
}

func (s *RobotTestSuite) TestStartSession() {
	robots, err := s.th.Service.Robot.List(s.ctx, "", "")
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(robots))

	areas, err := s.th.Service.Area.List(s.ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(areas))

	sess, err := s.th.Service.Robot.StartSession(s.ctx, entity.StartSessionArgs{
		RobotID:   robots[0].UID,
		AreaID:    areas[0].UID,
		RobotX:    3,
		RobotY:    4,
		StartedAt: time.Now(),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 3, sess.LastX)
	require.Equal(s.T(), 4, sess.LastY)
	require.NotEqual(s.T(), robots[0].Session[0].UID, sess.UID, "should have created a new cleaning session")
}

func (s *RobotTestSuite) TestUpdateSession() {
	robots, err := s.th.Service.Robot.List(s.ctx, "", "")
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(robots))

	areas, err := s.th.Service.Area.List(s.ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(areas))

	// Create a new session.
	sess, err := s.th.Service.Robot.StartSession(s.ctx, entity.StartSessionArgs{
		RobotID:   robots[0].UID,
		AreaID:    areas[0].UID,
		RobotX:    3,
		RobotY:    4,
		StartedAt: time.Now(),
	})
	require.NoError(s.T(), err)
	require.NotEqual(s.T(), robots[0].Session[0].UID, sess.UID, "should have created a new cleaning session")

	// Update session moving the robot until we reach 100% completion.
	center := robots[0].Size / 2
	robotX := 0
	robotY := 0
	reportedAt := time.Now()
	var secondsElapsed time.Duration

	var updSess *entity.CleaningSession
	for pass := 0; pass < sess.Area[0].PassesNeeded; pass++ {
		for _, sq := range sess.Area[0].Grid {
			robotX = sq.X + center
			robotY = sq.Y + center
			secondsElapsed++

			updSess, err = s.th.Service.Robot.UpdateSession(s.ctx, entity.UpdateSessionArgs{
				RobotID:    robots[0].UID,
				RobotX:     robotX,
				RobotY:     robotY,
				ReportedAt: reportedAt.Add(secondsElapsed * time.Second),
			})
			require.NoError(s.T(), err)
			require.Equal(s.T(), updSess.UID, sess.UID, "should have created a new cleaning session")
		}
	}
	require.Equal(s.T(), "100.00", updSess.Area[0].Completion(), "should have completed the entire cleaning area")

	// End session.
	endedAt := reportedAt.Add(secondsElapsed * time.Second)

	updSess, err = s.th.Service.Robot.EndSession(s.ctx, entity.UpdateSessionArgs{
		RobotID:    robots[0].UID,
		RobotX:     robotX,
		RobotY:     robotY,
		ReportedAt: endedAt,
	})
	require.GreaterOrEqual(s.T(), updSess.DurationSec, 24, "should have a duration of at least 24 seconds")

	// Fetch all history.
	history, err := s.th.Service.Robot.History(s.ctx, robots[0].UID, 10) // Fetch only the 10 latest cleaning sessions.
	require.GreaterOrEqual(s.T(), len(history.Session), 1, "should have at least one session")
	require.GreaterOrEqual(s.T(), len(history.Session[0].PositionHistory), 1, "should have at least one historical position")
	require.Equal(s.T(), endedAt.Unix(), history.Session[0].EndedAt.Unix(), "should have latest session ended at a predicatable time")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRobotTestSuite(t *testing.T) {
	ts := new(RobotTestSuite)
	ts.ctx = context.Background()
	ts.th = newTestHelper()
	suite.Run(t, ts)
}
