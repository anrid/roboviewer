// Package cerr helps us send more meaningful error responses
// to API clients.
package cerr

import (
	stderr "errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrServer is a general server error.
	ErrServer = stderr.New("server_error")

	// ErrValidationFailed means an payload (struct) failed validation.
	ErrValidationFailed = stderr.New("validation_failed")

	// ErrNotFound means that a resource could not be found.
	ErrNotFound = stderr.New("not_found")
)

// ErrorResponse is an error response.
type ErrorResponse struct {
	Ok    bool   `json:"ok"`
	Code  string `json:"code"`
	Error string `json:"error"`
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// PrintStackTrace prints a stack trace to stdout if the given
// error contains one.
func PrintStackTrace(err error) {
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			s := fmt.Sprintf("%+s:%d", f, f)
			if !strings.Contains(s, "testing.") && !strings.Contains(s, "runtime.") {
				println(s)
			}
		}
	}
}

// FromError converts the passed error to our custom Error type
// and also returns the related HTTP status code.
func FromError(err error) (e ErrorResponse, status int) {
	// Default error response and status code.
	e.Error = err.Error()
	e.Code = ErrServer.Error()
	status = http.StatusBadRequest

	c := errors.Cause(err)
	switch c {
	case ErrValidationFailed:
		e.Code = c.Error()
		status = http.StatusBadRequest
	case ErrNotFound:
		e.Code = c.Error()
		status = http.StatusNotFound
	default:
		log.Printf("unhandled error cause '%s'", c.Error())
	}
	return
}
