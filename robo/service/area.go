package service

import (
	"context"

	"github.com/anrid/roboviewer/robo/entity"
)

// AreaService holds all the route handlers (endpoints)
// related to robots.
type AreaService struct {
	r entity.AreaRepository
}

// NewAreaService creates a new robot controller instance.
func NewAreaService(r entity.AreaRepository) *AreaService {
	return &AreaService{r}
}

// List returns a list of all robots.
func (co *AreaService) List(ctx context.Context) ([]*entity.Area, error) {
	res, err := co.r.List(ctx)
	if err != nil {
		return nil, err
	}
	return res.Areas, nil
}
