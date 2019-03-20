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
	"github.com/labstack/gommon/color"
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
	var rwl = &responseWriterLogger{rw: rw}

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

// Filter a web request
func (r *logger) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	start := time.Now()

	webResponse := chain.Next(ctx, req, w)

	return &loggedResponse{
		result: webResponse,
		logCallback: func(rwl *responseWriterLogger) {
			var cp func(msg interface{}, styles ...string) string
			switch {
			case rwl.statusCode >= 200 && rwl.statusCode < 300:
				cp = color.Green
			case rwl.statusCode >= 300 && rwl.statusCode < 400:
				cp = color.Blue
			case rwl.statusCode >= 400 && rwl.statusCode < 500:
				cp = color.Yellow
			case rwl.statusCode == 0 || (rwl.statusCode >= 500 && rwl.statusCode < 600):
				cp = color.Red
			default:
				cp = color.Grey
			}

			extra := new(strings.Builder)

			switch r := webResponse.(type) {
			case *web.URLRedirectResponse:
				extra.WriteString("-> " + r.URL.String())

			case *web.RouteRedirectResponse:
				extra.WriteString("-> " + r.To)

			case *web.ServerErrorResponse:
				extra.WriteString(strings.Split(fmt.Sprintf(`Error: %s`, r.Error.Error()), "\n")[0])
			}

			var sizeStr string
			if rwl.length > 99999 {
				sizeStr = strconv.Itoa(rwl.length/1000) + "kb"
			} else {
				sizeStr = strconv.Itoa(rwl.length) + "b"
			}

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
		},
	}
}
