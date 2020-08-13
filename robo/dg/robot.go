package dg

import (
	"context"
	"strconv"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/dgraph-io/dgo/v2"
	"github.com/pkg/errors"
)

// RobotRepository ...
type RobotRepository struct {
	Repository
}

// NewRobotRepository creates a new repository.
func NewRobotRepository(c *dgo.Dgraph) *RobotRepository {
	return &RobotRepository{Repository: Repository{c}}
}

// List returns a list of robots together with their currently
// active cleaning session.
func (r *RobotRepository) List(ctx context.Context, a entity.ListRobotsArgs) (*entity.ListRobotsResult, error) {
	qb := NewQB(`
	query q($robotID: string, $name: string) {
		robots(func: type(Robot)) <FILTERS> {
			uid
			name
			size
			session @filter(eq(is_active, true)) (first: 1) (orderdesc: created_at) {
				uid
				name
				is_active
				started_at
				ended_at
				last_x
				last_y
				last_reported_at
				area {
					uid
					name
					size_x
					size_y
					passes_needed
					grid (orderasc: order) {
						uid
						x
						y
						size
						passes
						cleaned_at
						order
					}
				}
			}
		}
	}
	`)

	if a.RobotID != "" {
		qb.Filter(`uid($robotID)`)
	}
	if a.Name != "" {
		qb.Filter(`alloftext(name, $name)`)
	}
	query := qb.Query()
	// println(query)

	vars := map[string]string{
		"$robotID": a.RobotID,
		"$name":    a.Name,
	}
	resp, err := r.c.NewTxn().QueryWithVars(ctx, query, vars)
	if err != nil {
		return nil, err
	}
	// println(string(resp.Json))

	res := &entity.ListRobotsResult{}
	err = json.Unmarshal(resp.Json, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ByName returns one or more robots matching the
// name exactly.
func (r *RobotRepository) ByName(ctx context.Context, name string) (*entity.ListRobotsResult, error) {
	qb := NewQB(`
	query q($name: string) {
		robots(func: type(Robot), orderdesc: created_at) @filter(eq(name, $name)) {
			uid
			name
			is_cleaning
			size
			created_at
		}
	}
	`)
	query := qb.Query()
	// println(query)

	vars := map[string]string{
		"$name": name,
	}
	resp, err := r.c.NewTxn().QueryWithVars(ctx, query, vars)
	if err != nil {
		return nil, err
	}
	// println(string(resp.Json))

	res := &entity.ListRobotsResult{}
	err = json.Unmarshal(resp.Json, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetRobotAndArea returns a robot and an area by id.
func (r *RobotRepository) GetRobotAndArea(ctx context.Context, robotID, areaID string) (*entity.GetRobotAndAreaResult, error) {
	qb := NewQB(`
	query q($robotID: string, $areaID: string) {
		robots(func: type(Robot)) @filter(uid($robotID)) {
			uid
			name
			size
			created_at
			dgraph.type
			session @filter(eq(is_active, true)) (first: 1) (orderdesc: created_at) {
				uid
				name
				is_active
				started_at
				ended_at
				dgraph.type
			}
		}
		areas(func: type(Area)) @filter(uid($areaID)) {
			uid
			name
			size_x
			size_y
			passes_needed
			created_at
			dgraph.type
		}
	}
	`)
	query := qb.Query()

	vars := map[string]string{
		"$robotID": robotID,
		"$areaID":  areaID,
	}
	resp, err := r.c.NewTxn().QueryWithVars(ctx, query, vars)
	if err != nil {
		return nil, err
	}
	// println(string(resp.Json))

	res := &entity.GetRobotAndAreaResult{}
	err = json.Unmarshal(resp.Json, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// History returns all historial data for the given robot.
func (r *RobotRepository) History(ctx context.Context, robotID string, max int) (*entity.Robot, error) {
	qb := NewQB(`
	query q($robotID: string, $max: int) {
		robots(func: uid($robotID)) {
			uid
			name
			size
			session (first: $max) (orderdesc: created_at) {
				uid
				name
				is_active
				started_at
				ended_at
				last_x
				last_y
				last_reported_at
				duration_sec
				position_history (orderasc: passed_at) {
					x
					y
					passed_at	
				}
				area {
					uid
					name
					size_x
					size_y
					passes_needed
					grid (orderasc: order) {
						uid
						x
						y
						size
						passes
						cleaned_at
						order
					}
				}
			}
		}
	}
	`)
	query := qb.Query()
	// println(query)

	vars := map[string]string{
		"$robotID": robotID,
		"$max":     strconv.Itoa(max),
	}
	resp, err := r.c.NewTxn().QueryWithVars(ctx, query, vars)
	if err != nil {
		return nil, err
	}
	// println(string(resp.Json))

	res := &entity.ListRobotsResult{}
	err = json.Unmarshal(resp.Json, res)
	if err != nil {
		return nil, err
	}

	if len(res.Robots) != 1 {
		return nil, errors.Errorf("could not find robot %s", robotID)
	}
	return res.Robots[0], nil
}
