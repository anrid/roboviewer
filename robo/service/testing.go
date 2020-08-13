package service

import (
	"context"
	"sync"

	"github.com/anrid/roboviewer/robo/config"
	"github.com/anrid/roboviewer/robo/dg"
	"github.com/anrid/roboviewer/robo/entity"
)

var (
	once sync.Once
	th   *testHelper
)

type testHelper struct {
	Repository struct {
		Robot entity.RobotRepository
		Area  entity.AreaRepository
	}
	Service struct {
		Robot entity.RobotService
		Area  entity.AreaService
	}
}

// newTestHelper returns an initialized test helper.
func newTestHelper() *testHelper {
	once.Do(func() {
		c := config.GetConfig()

		th = &testHelper{}

		conn, _ := dg.Connect(c.DgraphURL)

		_ = dg.CreateSimpleTestData(context.Background(), conn)

		th.Repository.Robot = dg.NewRobotRepository(conn)
		th.Repository.Area = dg.NewAreaRepository(conn)

		th.Service.Robot = NewRobotService(th.Repository.Robot)
		th.Service.Area = NewAreaService(th.Repository.Area)
	})
	return th
}
