package controller

import (
	"github.com/anrid/roboviewer/robo/entity"
	"github.com/anrid/roboviewer/robo/pkg/httpserver"
	"github.com/labstack/echo/v4"
)

// AreaController holds all the route handlers (endpoints)
// related to areas.
type AreaController struct {
	svc entity.AreaService
}

// NewAreaController creates a new area controller instance.
func NewAreaController(svc entity.AreaService) *AreaController {
	return &AreaController{svc}
}

// ListAreas returns a list of all areas.
// @Summary     List all areas.
// @Description List all areas.
// @Accept      json
// @Produce     json
// @Success     200 {object} controller.ListAreasResponseV1
// @Failure     400 {object} cerr.ErrorResponse
// @Router      /v1/areas [get]
func (co *AreaController) ListAreas(c echo.Context) error {
	ctx := c.Request().Context()

	areas, err := co.svc.List(ctx)
	if err != nil {
		return httpserver.Fail(c, err)
	}

	return httpserver.Ok(c, ListAreasResponseV1{
		Ok:    true,
		Areas: areas,
	})
}

// ListAreasResponseV1 ...
type ListAreasResponseV1 struct {
	Ok    bool           `json:"ok"`
	Areas []*entity.Area `json:"areas"`
}

// SetupRoutes wires up the routes to the echo server.
func (co *AreaController) SetupRoutes(e *echo.Echo) {
	e.GET("/v1/areas", co.ListAreas)
}
