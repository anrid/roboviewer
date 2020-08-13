package controller

import (
	"strconv"

	"github.com/anrid/roboviewer/robo/entity"
	"github.com/anrid/roboviewer/robo/pkg/httpserver"
	"github.com/labstack/echo/v4"
)

// RobotController holds all the route handlers (endpoints)
// related to robots.
type RobotController struct {
	svc entity.RobotService
}

// NewRobotController creates a new robot controller instance.
func NewRobotController(svc entity.RobotService) *RobotController {
	return &RobotController{svc}
}

// List returns a list of all robots.
// @Summary     List all robots and their active cleaning session.
// @Description List all robots and their active cleaning session.
// @Accept      json
// @Produce     json
// @Param       robot_id query string false "Robot ID to filter on"
// @Param       name query string false "Robot name to filter on"
// @Success     200 {object} controller.ListRobotsResponseV1
// @Failure     400 {object} cerr.ErrorResponse
// @Router      /v1/robots [get]
func (co *RobotController) List(c echo.Context) error {
	ctx := c.Request().Context()

	robotID := c.QueryParam("robot_id")
	name := c.QueryParam("name")

	robots, err := co.svc.List(ctx, robotID, name)
	if err != nil {
		return httpserver.Fail(c, err)
	}

	return httpserver.Ok(c, ListRobotsResponseV1{
		Ok:     true,
		Robots: robots,
	})
}

// ListRobotsResponseV1 ...
type ListRobotsResponseV1 struct {
	Ok     bool            `json:"ok"`
	Robots []*entity.Robot `json:"robots"`
}

// History returns all historical data for a robot.
// @Summary     Get all historical cleaning sessions for a robot.
// @Description Get all historical cleaning sessions for a robot.
// @Accept      json
// @Produce     json
// @Param       robot_id path string true "Robot ID to show history for"
// @Param       max query integer false "Return only max latest number of cleaning sessions for robot (default: 10)"
// @Success     200 {object} controller.RobotHistoryResponseV1
// @Failure     400 {object} cerr.ErrorResponse
// @Router      /v1/robots/{robot_id}/history [get]
func (co *RobotController) History(c echo.Context) error {
	ctx := c.Request().Context()

	robotID := c.Param("robot_id")
	maxStr := c.QueryParam("max")
	max := 10
	if maxStr != "" {
		max, _ = strconv.Atoi(maxStr)
	}

	robot, err := co.svc.History(ctx, robotID, max)
	if err != nil {
		return httpserver.Fail(c, err)
	}

	return httpserver.Ok(c, RobotHistoryResponseV1{
		Ok:    true,
		Robot: robot,
	})
}

// RobotHistoryResponseV1 ...
type RobotHistoryResponseV1 struct {
	Ok    bool          `json:"ok"`
	Robot *entity.Robot `json:"robot"`
}

// SetupRoutes wires up the routes to the echo server.
func (co *RobotController) SetupRoutes(e *echo.Echo) {
	e.GET("/v1/robots", co.List)
	e.GET("/v1/robots/:robot_id/history", co.History)
}
