package requestlogger

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	logger struct {
		logger    flamingo.Logger
		responder *web.Responder
	}

	loggedResponse struct {
		result      web.Result
		logCallback func(rwl *responseWriterLogger)
	}

	responseWriterLogger struct {
		rw         http.ResponseWriter
		statusCode int
		length     int
	}
)

func (r *responseWriterLogger) Header() http.Header {
	return r.rw.Header()
}

func (r *responseWriterLogger) Write(b []byte) (int, error) {
	length, err := r.rw.Write(b)
	r.length += length
	return length, err
}

func (r *responseWriterLogger) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.rw.WriteHeader(statusCode)
}

// Apply logger to request
func (l *loggedResponse) Apply(ctx context.Context, rw http.ResponseWriter) error {
	var err error
	var rwl = &responseWriterLogger{rw: rw, statusCode: http.StatusOK}

	defer l.logCallback(rwl)

	if l.result != nil {
		err = l.result.Apply(ctx, rwl)
	}

	return err
}

func (r *logger) Inject(flogger flamingo.Logger, responder *web.Responder) {
	r.logger = flogger
	r.responder = responder
}

func humanBytes(bc int) string {
	if bc > 99999 {
		return strconv.Itoa(bc/1000) + "kb"
	}
	return strconv.Itoa(bc) + "b"
}

const (
	green  = "32"
	blue   = "34"
	yellow = "33"
	red    = "31"
	grey   = "90"
)

func colored(color string) func(msg string) string {
	return func(msg string) string {
		return fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, msg)
	}
}

func statusCodeColor(code int) func(msg string) string {
	switch {
	case code >= 200 && code < 300:
		return colored(green)
	case code >= 300 && code < 400:
		return colored(blue)
	case code >= 400 && code < 500:
		return colored(yellow)
	case code == 0 || (code >= 500 && code < 600):
		return colored(red)
	default:
		return colored(grey)
	}
}

// Filter a web request
func (r *logger) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	start := time.Now()

	webResponse := chain.Next(ctx, req, w)

	logCallbackFunc := func(rwl *responseWriterLogger) {
		cp := statusCodeColor(rwl.statusCode)
		extra := new(strings.Builder)

		switch r := webResponse.(type) {
		case *web.URLRedirectResponse:
			extra.WriteString("-> " + r.URL.String())

		case *web.RouteRedirectResponse:
			extra.WriteString("-> " + r.To)

		case *web.ServerErrorResponse:
			extra.WriteString(strings.Split(fmt.Sprintf(`Error: %s`, r.Error), "\n")[0])
		}

		sizeStr := humanBytes(rwl.length)

		duration := time.Since(start)

		l := r.logger.
			WithContext(ctx).
			WithFields(
				map[flamingo.LogKey]interface{}{
					flamingo.LogKeyAccesslog:    1,
					flamingo.LogKeyResponseCode: rwl.statusCode,
					flamingo.LogKeyResponseTime: duration,
					flamingo.LogKeyReferer:      req.Request().Referer(),
					flamingo.LogKeyClientIP:     strings.Join(req.RemoteAddress(), ", "),
					flamingo.LogKeyBusinessID:   req.Request().Header.Get("X-Business-ID"),
				},
			)

		var extraStr string
		if extra.Len() > 0 {
			extraStr = " (" + extra.String() + ")"
		}

		l.Info(
			fmt.Sprintf(
				cp("%s %s %d: %s in %s%s"),
				req.Request().Method,
				req.Request().RequestURI,
				rwl.statusCode,
				sizeStr,
				duration,
				extraStr,
			),
		)
	}
	return &loggedResponse{
		result:      webResponse,
		logCallback: logCallbackFunc,
	}
}
