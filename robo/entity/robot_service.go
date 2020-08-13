package entity

import (
	"context"
	"time"
)

// RobotService holds various use cases related to robots.
type RobotService interface {
	List(ctx context.Context, id, name string) ([]*Robot, error)
	StartSession(ctx context.Context, a StartSessionArgs) (*CleaningSession, error)
	UpdateSession(ctx context.Context, a UpdateSessionArgs) (*CleaningSession, error)
	EndSession(ctx context.Context, a UpdateSessionArgs) (*CleaningSession, error)
	History(ctx context.Context, robotID string, max int) (*Robot, error)
}

// StartSessionArgs are passed to RobotService.StartSession.
type StartSessionArgs struct {
	RobotID   string    // RobotID of the robot to do the cleaning.
	AreaID    string    // AreaID of the area to clean.
	RobotX    int       // Robot's initial X coordinate (optional).
	RobotY    int       // Robot's initial Y coordinate (optional).
	StartedAt time.Time // When the session started according to the robot.
}

// UpdateSessionArgs are passed to RobotService.StartSession.
type UpdateSessionArgs struct {
	RobotID    string    // RobotID of the robot to do the cleaning.
	RobotX     int       // Robot's current X coordinate.
	RobotY     int       // Robot's current Y coordinate.
	ReportedAt time.Time // When this position was reported according to the robot.
	EndSession bool      // Close the session.
}
