package controller

import (
	"sync"

	"github.com/anrid/roboviewer/robo/pkg/testserver"
)

var (
	once sync.Once
	ts   *testserver.TS
)

// Setup sets up a test environment when testing handlers.
func setupTests() *testserver.TS {
	once.Do(func() {
		ts = testserver.Get()
		NewRobotController(ts.Service.Robot).SetupRoutes(ts.Server.Echo)
		NewAreaController(ts.Service.Area).SetupRoutes(ts.Server.Echo)
	})
	return ts
}
