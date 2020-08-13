package controller

import (
	"net/http"
	"testing"

	"github.com/anrid/roboviewer/robo/pkg/httpserver"
	"github.com/stretchr/testify/require"
)

func TestListAreas(t *testing.T) {
	ts := setupTests()

	out := &ListAreasResponseV1{}

	status, _ := httpserver.Call(http.MethodGet, "/v1/areas", ts.Server, nil, out)
	require.Equal(t, http.StatusOK, status, "should succeed")
	require.Equal(t, 2, len(out.Areas), "should list 2 areas")
}
