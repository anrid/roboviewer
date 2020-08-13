package entity

import (
	"fmt"
	"time"
)

// CleaningSession is a robot cleaning session.
type CleaningSession struct {
	Name            string          `json:"name,omitempty"` // Optional.
	Area            []*CleaningArea `json:"area,omitempty"`
	IsActive        bool            `json:"is_active,omitempty"`
	StartedAt       *time.Time      `json:"started_at,omitempty"`
	EndedAt         *time.Time      `json:"ended_at,omitempty"`
	LastX           int             `json:"last_x,omitempty"`
	LastY           int             `json:"last_y,omitempty"`
	LastReportedAt  *time.Time      `json:"last_reported_at,omitempty"`
	PositionHistory []*Position     `json:"position_history,omitempty"`
	DurationSec     int             `json:"duration_sec,omitempty"`
	Common
}

// NewCleaningSession creates a new cleaning session for a
// robot in a given area.
func NewCleaningSession(r *Robot, a *Area, name string) *CleaningSession {
	if name == "" {
		// Let's try to create informative default names
		// for our cleaning sessions.
		name = fmt.Sprintf("Cleaning session: robot %s in area %s on %s", r.Name, a.Name, now().Format("2 Jan 2006 15:04"))
	}
	cs := &CleaningSession{
		Name: name,
		Area: []*CleaningArea{
			NewCleaningArea(a, r),
		},
		IsActive:  true,
		StartedAt: now(),
		Common: Common{
			UID:       "_:" + CleaningSessionUID,
			DType:     []string{"CleaningSession"},
			CreatedAt: now(),
		},
	}
	return cs
}

// End ends this cleaning session if it was active.
func (cs *CleaningSession) End(endedAt time.Time) {
	if cs.IsActive {
		// Mark session as inactive.
		cs.EndedAt = &endedAt
		cs.IsActive = false
		if cs.StartedAt != nil {
			// Calculate session duration.
			cs.DurationSec = int(cs.EndedAt.Sub(*cs.StartedAt).Seconds())
		}
	}
}
