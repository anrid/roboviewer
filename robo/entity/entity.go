package entity

import (
	"encoding/json"
	"time"
)

const (
	// RobotUID ...
	RobotUID = "r"
	// AreaUID ...
	AreaUID = "a"
	// CleaningSessionUID ...
	CleaningSessionUID = "cs"
	// CleaningAreaUID ...
	CleaningAreaUID = "ca"
	// SquareUID ...
	SquareUID = "sq"
	// PositionUID ...
	PositionUID = "p"
)

// Common contains common mandatory fields used
// in all our entities.
type Common struct {
	UID       string     `json:"uid,omitempty"`
	DType     []string   `json:"dgraph.type,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func now() *time.Time {
	t := time.Now()
	return &t
}

// Dump marshals the given object to JSON and
// pretty prints it.
func Dump(v interface{}) {
	d, _ := json.MarshalIndent(v, "", "  ")
	println("DUMP:")
	println(string(d))
}
