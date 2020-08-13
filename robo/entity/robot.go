package entity

// Robot is a vacuum cleaning robot.
type Robot struct {
	// Each robot should have a name to make identification easier and reports nicer.
	Name string `json:"name,omitempty"`

	IsCleaning bool               `json:"is_cleaning,omitempty"`
	Session    []*CleaningSession `json:"session,omitempty"`

	// The diameter of the robot in millimeters. We assume all robots are have a circle shape.
	Size int `json:"size,omitempty"`

	Common
}

// NewCleaningSession is a helper method to quickly create a
// new cleaning session for a robot.
func (r *Robot) NewCleaningSession(a *Area) *CleaningSession {
	newSess := NewCleaningSession(r, a, "")
	r.Session = []*CleaningSession{newSess}
	return newSess
}

// NewRobot creates a new robot with a name and size.
func NewRobot(name string, size int) *Robot {
	return &Robot{
		Name: name,
		Size: size,
		Common: Common{
			UID:       "_:" + RobotUID,
			DType:     []string{"Robot"},
			CreatedAt: now(),
		},
	}
}
