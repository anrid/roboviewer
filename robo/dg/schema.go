package dg

import (
	"context"
	"log"
	"time"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
)

// CreateSchema ...
func CreateSchema(ctx context.Context, c *dgo.Dgraph) {
	t := time.Now()

	op := &api.Operation{}
	op.Schema = `
		# String fields
		name: string @index(fulltext) .

		# Int fields
		size: int .
		x: int . 
		y: int . 
		last_x: int . 
		last_y: int . 
		size_x: int .
		size_y: int .
		passes_needed: int .
		passes: int .
		order: int @index(int) .
		duration_sec: int .

		# Date fields
		started_at: dateTime @index(hour) .
		ended_at: dateTime @index(hour) .
		cleaned_at: dateTime @index(hour) .
		created_at: dateTime @index(hour) .
		passed_at: dateTime @index(hour) .
		last_reported_at: dateTime .

		# Boolean fields
		is_active: bool @index(bool) .

		# Edges (joinable)
		robot: [uid] @reverse .
		grid: [uid] @reverse .
		session: [uid] @reverse . 
		area: [uid] @reverse . 
		position_history: [uid] .

		type Robot {
			name
			session
			size
			created_at
		}

		type CleaningSession {
			name
			area
			is_active
			started_at
			ended_at
			created_at
			last_x
			last_y
			last_reported_at
			position_history
			duration_sec
		}

		type Area {
			name
			size_x
			size_y
			passes_needed
			created_at
		}
		
		type CleaningArea {
			name
			size_x
			size_y
			grid
			passes_needed
			created_at
		}

		type Square {
			x
			y
			size
			passes
			cleaned_at
			order
		}

		type Position {
			x
			y
			passed_at
			order
		}
	`

	// Setup schema.
	if err := c.Alter(ctx, op); err != nil {
		log.Fatalf("could not alter schema: %s", err.Error())
	}
	log.Printf("created database schema in %s", time.Since(t).String())
}

// CreateSimpleTestData ...
func CreateSimpleTestData(ctx context.Context, c *dgo.Dgraph) map[string]string {
	timer := time.Now()

	var robo1 *entity.Robot
	var robo2 *entity.Robot

	robo1 = entity.NewRobot("Test - Johnny 5", 500)
	robo2 = entity.NewRobot("Test - ED 209", 1000)

	area1 := entity.NewArea("Tiny Room 1", 1000, 2000, 3)
	area2 := entity.NewArea("Tiny Room 2", 2000, 2000, 2)

	robo1.NewCleaningSession(area1)
	robo2.NewCleaningSession(area2)

	// Find existing test robots and use their UIDs.
	{
		repo := NewRobotRepository(c)

		for _, r := range []*entity.Robot{robo1, robo2} {
			list, err := repo.ByName(ctx, r.Name)
			if err != nil {
				panic(err)
			}
			if len(list.Robots) > 0 {
				r.UID = list.Robots[0].UID
				log.Printf("found existing robot: UID=%s name=%s", r.UID, r.Name)
			}
		}
	}

	// Find existing test areas and use their UIDs.
	{
		repo := NewAreaRepository(c)

		list, err := repo.List(ctx)
		if err != nil {
			panic(err)
		}

		for _, a := range list.Areas {
			if a.Name == area1.Name {
				area1.UID = a.UID
				log.Printf("found existing area: UID=%s name=%s", area1.UID, area1.Name)
			}
			if a.Name == area2.Name {
				area2.UID = a.UID
				log.Printf("found existing area: UID=%s name=%s", area2.UID, area2.Name)
			}
		}
	}

	// spew.Dump(robo1)
	// spew.Dump(robo2)

	pk1, err := Store(ctx, c, robo1)
	if err != nil {
		panic(err)
	}
	pk2, err := Store(ctx, c, robo2)
	if err != nil {
		panic(err)
	}
	pk3, err := Store(ctx, c, area1)
	if err != nil {
		panic(err)
	}
	pk4, err := Store(ctx, c, area2)
	if err != nil {
		panic(err)
	}

	uids := map[string]string{
		"r1": pk1[entity.RobotUID],
		"r2": pk2[entity.RobotUID],
		"a1": pk3[entity.AreaUID],
		"a2": pk4[entity.AreaUID],
		"s1": pk1[entity.CleaningSessionUID],
		"s2": pk2[entity.CleaningSessionUID],
	}

	spew.Dump(uids)

	log.Printf("created simple test data in %s", time.Since(timer))
	return uids
}
