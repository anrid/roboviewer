package entity

import (
	"context"
)

// AreaService holds various use cases related to areas.
type AreaService interface {
	List(ctx context.Context) ([]*Area, error)
}
