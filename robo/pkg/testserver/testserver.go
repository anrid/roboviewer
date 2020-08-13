package testserver

import (
	"context"
	"sync"

	"github.com/anrid/roboviewer/robo/config"
	"github.com/anrid/roboviewer/robo/dg"
	"github.com/anrid/roboviewer/robo/entity"
	"github.com/anrid/roboviewer/robo/pkg/httpserver"
	"github.com/anrid/roboviewer/robo/service"
)

var (
	once sync.Once
	ts   *TS
)

// TS is a test server.
type TS struct {
	Server     *httpserver.Server
	Repository struct {
		Robot entity.RobotRepository
		Area  entity.AreaRepository
	}
	Service struct {
		Robot entity.RobotService
		Area  entity.AreaService
	}
}

// Get creates a test server used to test handlers.
func Get() *TS {
	once.Do(func() {
		c := config.GetConfig()

		ts = &TS{
			Server: httpserver.NewServer(),
		}

		conn, _ := dg.Connect(c.DgraphURL)

		_ = dg.CreateSimpleTestData(context.Background(), conn)

		ts.Repository.Robot = dg.NewRobotRepository(conn)
		ts.Repository.Area = dg.NewAreaRepository(conn)

		ts.Service.Robot = service.NewRobotService(ts.Repository.Robot)
		ts.Service.Area = service.NewAreaService(ts.Repository.Area)
	})
	return ts
}
