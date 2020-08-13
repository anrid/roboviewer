package entity

import "context"

// AreaRepository defines data layer functionality related to areas.
type AreaRepository interface {
	List(ctx context.Context) (*ListAreasResult, error)
	Repository
}

// ListAreasResult is a list of areas together with their active
// cleaning session.
type ListAreasResult struct {
	Areas []*Area `json:"areas"`
}
