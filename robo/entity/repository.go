package entity

import (
	"context"
)

// Repository contains general persistance layer operations.
type Repository interface {
	Save(ctx context.Context, object interface{}) (map[string]string, error)
}
