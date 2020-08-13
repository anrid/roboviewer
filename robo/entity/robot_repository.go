package entity

import "context"

// RobotRepository defines data layer functionality related to robots.
type RobotRepository interface {
	List(ctx context.Context, a ListRobotsArgs) (*ListRobotsResult, error)
	GetRobotAndArea(ctx context.Context, robotID, areaID string) (*GetRobotAndAreaResult, error)
	History(ctx context.Context, robotID string, max int) (*Robot, error)
	Repository
}

// ListRobotsArgs are the args we pass to RobotRepository.List().
type ListRobotsArgs struct {
	RobotID string
	Name    string
}

// ListRobotsResult is a list of robots together with their active
// cleaning session.
type ListRobotsResult struct {
	Robots []*Robot `json:"robots"`
}

// GetRobotAndAreaResult is a list of robots together with their active
// cleaning session.
type GetRobotAndAreaResult struct {
	Robots []*Robot `json:"robots"`
	Areas  []*Area  `json:"areas"`
}
