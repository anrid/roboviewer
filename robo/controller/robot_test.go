package controller

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/anrid/roboviewer/robo/pkg/httpserver"
	"github.com/stretchr/testify/require"
)

func TestRobots(t *testing.T) {
	ts := setupTests()

	var robots []*entity.Robot
	// List robots.
	{
		out := &ListRobotsResponseV1{}

		status, _ := httpserver.Call(http.MethodGet, "/v1/robots", ts.Server, nil, out)
		require.Equal(t, http.StatusOK, status, "should succeed")
		require.Equal(t, 2, len(out.Robots), "should list 2 robots")

		robots = out.Robots
	}

	// Get history data.
	{
		out := &RobotHistoryResponseV1{}

		status, _ := httpserver.Call(http.MethodGet, fmt.Sprintf("/v1/robots/%s/history", robots[0].UID), ts.Server, nil, out)
		require.Equal(t, http.StatusOK, status, "should succeed")
		require.NotEmpty(t, out.Robot, "should get 1 robot")
	}
}
