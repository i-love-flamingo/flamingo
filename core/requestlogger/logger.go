package requestlogger

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/router"
	"github.com/labstack/gommon/color"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	contextKey string
	logger     struct{
		Logger flamingo.Logger `inject:""`
	}
)

const contextTime contextKey = "time"

// Notify is called on events
func (r *logger) Notify(e event.Event) {
	switch e := e.(type) {
	case *router.OnRequestEvent:
		r.onRequest(e)

	case *router.OnFinishEvent:
		r.onFinish(e)
	}
}

// onRequest assigns the current time to the request-context
func (r *logger) onRequest(event *router.OnRequestEvent) {
	event.Request = event.Request.WithContext(context.WithValue(event.Request.Context(), contextTime, time.Now()))
}

// onFinish logs the request to stdout via log.Printf
func (r *logger) onFinish(event *router.OnFinishEvent) {
	var duration time.Duration

	response, ok := event.ResponseWriter.(*router.VerboseResponseWriter)
	if !ok {
		return
	}

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
		cp = color.Grey
	}

	var extra string

	if response.Header().Get("Location") != "" {
		extra += " -> " + response.Header().Get("Location")
	}

	if event.Error != nil {
		extra += strings.Split(fmt.Sprintf(` | Error: %s`, event.Error), "\n")[0]
	}

	r.Logger.Printf(
		cp("%03d | %-8s | % 15s | % 6d byte | %s%s"),
		response.Status,
		event.Request.Method,
		duration,
		response.Size,
		event.Request.RequestURI,
		extra,
	)
}
