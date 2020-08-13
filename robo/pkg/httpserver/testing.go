package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

// Call calls an endpoint on given server.
func Call(method, url string, serv *Server, in, out interface{}) (status int, body string) {
	var data io.Reader

	if in != nil {
		d, err := json.Marshal(in)
		if err != nil {
			panic("could not marshal payload to JSON")
		}
		data = bytes.NewReader(d)
	}

	req := httptest.NewRequest(method, url, data)

	// Force JSON Content-Type.
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	w := httptest.NewRecorder()
	serv.Echo.ServeHTTP(w, req)
	resp := w.Result()
	bytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 300 && out != nil {
		_ = json.Unmarshal(bytes, out)
	}
	return resp.StatusCode, string(bytes)
}
