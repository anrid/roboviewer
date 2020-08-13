package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anrid/roboviewer/robo/pkg/cerr"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

// Server is a simple HTTP server based on echo.
type Server struct {
	Echo *echo.Echo
}

// We create a custom validator so that we can use
// validator.v9 with Echo.
type customValidator struct {
	v *validator.Validate
}

func (c *customValidator) Validate(d interface{}) error {
	return c.v.Struct(d)
}

// NewServer creates a new server instance.
func NewServer() *Server {
	e := echo.New()

	// Setup custom validator.
	e.Validator = &customValidator{validator.New()}

	// Middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Basic routes.
	e.GET("/", rootHandler)
	e.POST("/test/echo", echoHandler)

	// Not found handler.
	e.GET("/*", func(c echo.Context) error {
		msg := fmt.Sprintf("Route not found: %s", c.Request().URL)
		return c.String(http.StatusNotFound, msg)
	})

	return &Server{e}
}

// Ok is called when a handler returns a successful response.
func Ok(c echo.Context, payload interface{}) error {
	return c.JSON(http.StatusOK, payload)
}

// Fail is called when a handler returns an error.
func Fail(c echo.Context, err error) error {
	// Let the world know.
	log.Printf("api error: %s", err.Error())

	cerr.PrintStackTrace(err)

	resp, code := cerr.FromError(err)
	return c.JSON(code, resp)
}

// Bind binds the incoming request body to a struct and validates.
func Bind(c echo.Context, payload interface{}) (err error) {
	err = c.Bind(payload)
	if err != nil {
		return
	}
	return c.Validate(payload)
}

func rootHandler(c echo.Context) error {
	return Ok(c, echo.Map{
		"ok":        true,
		"message":   "All is well in the world!",
		"timestamp": time.Now().UnixNano(),
	})
}

func echoHandler(c echo.Context) error {
	r := &EchoRequestV1{}
	err := Bind(c, r)
	if err != nil {
		return Fail(c, err)
	}

	return Ok(c, EchoResponseV1{
		Ok:        true,
		Message:   r.Message,
		Timestamp: time.Now().UnixNano(),
	})
}

// EchoRequestV1 ...
type EchoRequestV1 struct {
	Message string `json:"message" validate:"required,gte=1"`
}

// EchoResponseV1 ...
type EchoResponseV1 struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
