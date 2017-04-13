package requestlogger

import (
	"context"
	"flamingo/framework/router"
	"fmt"
	"log"
	"time"

	"github.com/labstack/gommon/color"
)

type contextKey string

const contextTime contextKey = "time"

// Logger logs requests to stdout
type Logger struct{}

// Events returns a list of subscribed events
func (r *Logger) Events() []interface{} {
	return []interface{}{
		router.REQUEST,
		router.FINISH,
	}
}

// OnRequest assigns the current time to the request-context
func (r *Logger) OnRequest(event *router.OnRequestEvent) {
	event.Request = event.Request.WithContext(context.WithValue(event.Request.Context(), contextTime, time.Now()))
}

// OnFinish logs the request to stdout via log.Printf
func (r *Logger) OnFinish(event *router.OnFinishEvent) {
	var duration time.Duration
	var response, _ = event.ResponseWriter.(*router.VerboseResponseWriter)

	if start := event.Request.Context().Value(contextTime); start != nil {
		duration = time.Since(start.(time.Time))
	}

	var cp func(msg interface{}, styles ...string) string
	switch {
	case response.Status >= 200 && response.Status < 300:
		cp = color.Green
	case response.Status >= 300 && response.Status < 400:
		cp = color.Blue
	case response.Status >= 400 && response.Status < 500:
		cp = color.Yellow
	case response.Status >= 500 && response.Status < 600:
		cp = color.Red
	default:
		cp = color.Black
	}

	var extra string

	if response.Header().Get("Location") != "" {
		extra += " -> " + response.Header().Get("Location")
	}

	if event.Error != nil {
		extra += fmt.Sprintf(` | Error: %s`, event.Error)
	}

	log.Printf(
		cp("%03d | %-8s | % 15s | % 6d byte | %s%s"),
		response.Status,
		event.Request.Method,
		duration,
		response.Size,
		event.Request.RequestURI,
		extra,
	)
}
