package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/gommon/color"
)

type (
	// ResponseWriterLogger defines the function ResponseWriter uses for debugging/logging
	ResponseWriterLogger interface {
		Printf(string, ...interface{})
	}

	// ResponseWriter shadows http.ResponseWriter and tracks written bytes and result status for logging.
	ResponseWriter struct {
		http.ResponseWriter
		status int
		size   int
	}
)

// Write calls http.ResponseWriter.Write and records the written bytes.
func (response *ResponseWriter) Write(data []byte) (int, error) {
	l, e := response.ResponseWriter.Write(data)
	response.size += l
	return l, e
}

// WriteHeader calls http.ResponseWriter.WriteHeader and records the status code.
func (response *ResponseWriter) WriteHeader(h int) {
	response.status = h
	response.ResponseWriter.WriteHeader(h)
}

// Log Requests in a very dirty way.
func (response *ResponseWriter) Log(logger ResponseWriterLogger, duration time.Duration, req *http.Request, err interface{}) {
	var cp func(msg interface{}, styles ...string) string
	switch {
	case response.status >= 200 && response.status < 300:
		cp = color.Green
	case response.status >= 300 && response.status < 400:
		cp = color.Blue
	case response.status >= 400 && response.status < 500:
		cp = color.Yellow
	case response.status >= 500 && response.status < 600:
		cp = color.Red
	default:
		cp = color.Black
	}

	var extra string

	if response.Header().Get("Location") != "" {
		extra += " -> " + response.Header().Get("Location")
	}

	if err != nil {
		extra += fmt.Sprintf(` | Error: %s`, err)
	}

	logger.Printf(
		cp("%03d | %-8s | % 15s | % 6d byte | %s%s"),
		response.status,
		req.Method,
		duration,
		response.size,
		req.RequestURI,
		extra,
	)
}
