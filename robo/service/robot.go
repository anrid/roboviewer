package service

import (
	"context"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/anrid/roboviewer/robo/pkg/cerr"
	"github.com/pkg/errors"
)

// RobotService holds all the route handlers (endpoints)
// related to robots.
type RobotService struct {
	r entity.RobotRepository
}

// NewRobotService creates a new robot controller instance.
func NewRobotService(r entity.RobotRepository) *RobotService {
	return &RobotService{r}
}

// List returns a list of all robots.
func (co *RobotService) List(ctx context.Context, id, name string) ([]*entity.Robot, error) {
	res, err := co.r.List(ctx, entity.ListRobotsArgs{
		RobotID: id,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}
	return res.Robots, nil
}

// StartSession starts a new cleaning session for the given
// robot and area.
func (co *RobotService) StartSession(ctx context.Context, a entity.StartSessionArgs) (*entity.CleaningSession, error) {
	if a.StartedAt.IsZero() {
		return nil, errors.Wrapf(cerr.ErrValidationFailed, "did not get a valid StartedAt value: %s", a.StartedAt)
	}

	res, err := co.r.GetRobotAndArea(ctx, a.RobotID, a.AreaID)
	if err != nil {
		return nil, err
	}
	if len(res.Robots) == 0 {
		return nil, errors.Wrapf(cerr.ErrNotFound, "could not find robot with id %s", a.RobotID)
	}
	if len(res.Areas) == 0 {
		return nil, errors.Wrapf(cerr.ErrNotFound, "could not find area with id %s", a.AreaID)
	}

	robot := res.Robots[0]
	area := res.Areas[0]

	if len(robot.Session) > 0 {
		// End the ongoing session.
		prevSess := robot.Session[0]
		prevSess.End(a.StartedAt)
		_, err := co.r.Save(ctx, prevSess)
		if err != nil {
			return nil, errors.Wrap(err, "could not persist previous session")
		}
		robot.Session = nil
	}

	newSess := robot.NewCleaningSession(area)
	newSess.LastX = a.RobotX
	newSess.LastY = a.RobotY
	newSess.StartedAt = &a.StartedAt
	newSess.PositionHistory = []*entity.Position{entity.NewPosition(a.RobotX, a.RobotY, a.StartedAt)}

	uids, err := co.r.Save(ctx, robot)
	if err != nil {
		return nil, errors.Wrap(err, "could not persist new session")
	}

	// Assign new primary key to object.
	// TODO: This probably needs to be improved by
	// making it automatic, perhaps performed by the
	// dg.Repository.Save().
	newSess.UID = uids[entity.CleaningSessionUID]

	return newSess, nil
}

// UpdateSession updates a robot's current cleaning
// session, called every time a robot moves.
// It also allow us to close the session.
func (co *RobotService) UpdateSession(ctx context.Context, a entity.UpdateSessionArgs) (*entity.CleaningSession, error) {
	if a.ReportedAt.IsZero() {
		return nil, errors.Wrapf(cerr.ErrValidationFailed, "did not get a valid ReportedAt value: %s", a.ReportedAt)
	}

	res, err := co.r.List(ctx, entity.ListRobotsArgs{
		RobotID: a.RobotID,
	})
	if err != nil {
		return nil, err
	}
	if len(res.Robots) == 0 {
		return nil, errors.Wrapf(cerr.ErrNotFound, "could not find robot with id %s", a.RobotID)
	}

	robot := res.Robots[0]

	if len(robot.Session) == 0 {
		return nil, errors.Wrapf(cerr.ErrNotFound, "could not find any sessions for robot %s id %s", robot.Name, robot.UID)
	}
	if !robot.Session[0].IsActive {
		return nil, errors.Wrapf(cerr.ErrNotFound, "could not find an active session for robot %s id %s", robot.Name, robot.UID)
	}

	sess := robot.Session[0]

	sess.LastX = a.RobotX
	sess.LastY = a.RobotY
	sess.LastReportedAt = &a.ReportedAt
	sess.Area[0].SetVisited(a.RobotX, a.RobotY)
	sess.PositionHistory = []*entity.Position{entity.NewPosition(a.RobotX, a.RobotY, a.ReportedAt)}

	if a.EndSession {
		sess.End(a.ReportedAt)
	}

	_, err = co.r.Save(ctx, sess)
	if err != nil {
		return nil, errors.Wrap(err, "could not persist current session")
	}

	return sess, nil
}

// EndSession ends an active session for a given robot.
func (co *RobotService) EndSession(ctx context.Context, a entity.UpdateSessionArgs) (*entity.CleaningSession, error) {
	a.EndSession = true
	return co.UpdateSession(ctx, a)
}

// History gets all cleaning session and position history for a robot.
func (co *RobotService) History(ctx context.Context, robotID string, max int) (*entity.Robot, error) {
	return co.r.History(ctx, robotID, max)
}
