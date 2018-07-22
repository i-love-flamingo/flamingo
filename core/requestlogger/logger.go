package requestlogger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"github.com/labstack/gommon/color"
)

type (
	logger struct {
		logger flamingo.Logger
	}

	loggedResponse struct {
		web.Response
		logCallback func()
	}
)

// Apply logger to request
func (l *loggedResponse) Apply(ctx context.Context, rw http.ResponseWriter) {
	if l.Response != nil {
		l.Response.Apply(ctx, rw)
	}

	l.logCallback()
}

func (r *logger) Inject(flogger flamingo.Logger) {
	r.logger = flogger
}

// Filter a web request
func (r *logger) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	start := time.Now()

	webResponse := chain.Next(ctx, req, w)

	if webResponse == nil {
		webResponse = &web.ErrorResponse{Error: errors.New("nil response"), Response: &web.BasicResponse{}}
	}

	return &loggedResponse{
		Response: webResponse,
		logCallback: func() {
			duration := time.Since(start)

			var cp func(msg interface{}, styles ...string) string
			switch {
			case webResponse.GetStatus() >= 200 && webResponse.GetStatus() < 300:
				cp = color.Green
			case webResponse.GetStatus() >= 300 && webResponse.GetStatus() < 400:
				cp = color.Blue
			case webResponse.GetStatus() >= 400 && webResponse.GetStatus() < 500:
				cp = color.Yellow
			case webResponse.GetStatus() >= 500 && webResponse.GetStatus() < 600:
				cp = color.Red
			default:
				cp = color.Grey
			}

			var extra string

			if rr, ok := webResponse.(*web.RedirectResponse); ok {
				extra += " -> " + rr.Location
			}

			if r, ok := webResponse.(web.ErrorResponse); ok && r.Error != nil {
				extra += strings.Split(fmt.Sprintf(` | Error: %s`, r.Error.Error()), "\n")[0]
			}

			var sizeStr string
			if webResponse.GetContentLength() > 99999 {
				sizeStr = strconv.Itoa(webResponse.GetContentLength()/1000) + "kb"
			} else {
				sizeStr = strconv.Itoa(webResponse.GetContentLength()) + "b"
			}

			l := r.logger.
				WithContext(ctx).
				WithFields(
					map[flamingo.LogKey]interface{}{
						flamingo.LogKeyAccesslog:    1,
						flamingo.LogKeyResponseCode: webResponse.GetStatus(),
						flamingo.LogKeyResponseTime: duration,
					},
				)
			defer l.Flush()

			l.Info(
				fmt.Sprintf(
					cp("%03d | %-8s | % 15s | % 7s | %s%s"),
					webResponse.GetStatus(),
					req.Request().Method,
					duration,
					sizeStr,
					req.Request().RequestURI,
					extra,
				),
			)
		},
	}
}
