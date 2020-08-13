package dg

import (
	"context"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/dgraph-io/dgo/v2"
)

// AreaRepository ...
type AreaRepository struct {
	Repository
}

// NewAreaRepository creates a new repository.
func NewAreaRepository(c *dgo.Dgraph) *AreaRepository {
	return &AreaRepository{Repository: Repository{c}}
}

// List returns a list of areas.
func (r *AreaRepository) List(ctx context.Context) (*entity.ListAreasResult, error) {
	qb := NewQB(`
	query {
		areas(func: type(Area), first: 100, orderdesc: created_at) {
			uid
			name
			size_x
			size_y
			passes_needed
		}
	}
	`)

	query := qb.Query()
	// println(query)

	resp, err := r.c.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	// println(string(resp.Json))

	res := &entity.ListAreasResult{}
	err = json.Unmarshal(resp.Json, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
