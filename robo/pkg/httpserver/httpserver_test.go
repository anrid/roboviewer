package httpserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	serv := NewServer()
	{
		status, body := Call(http.MethodGet, "/", serv, nil, nil)
		require.Equal(t, http.StatusOK, status, "should get a response")
		require.Contains(t, body, `"message":"All is well in the world!"`, "should return 'All is well' message")
	}
	{
		in := &EchoRequestV1{Message: "Whut?!"}
		out := &EchoResponseV1{}
		status, _ := Call(http.MethodPost, "/test/echo", serv, in, out)
		require.Equal(t, http.StatusOK, status, "should get a response")
		require.Equal(t, in.Message, out.Message, "should get the sent message back as an echo")
	}
}
