package requestlogger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/labstack/gommon/color"
	"github.com/pkg/errors"
)

type (
	logger struct {
		logger    flamingo.Logger
		responder *web.Responder
	}

	loggedResponse struct {
		web.Result
		logCallback func()
	}
)

// Apply logger to request
func (l *loggedResponse) Apply(ctx context.Context, rw http.ResponseWriter) error {
	if l.Result != nil {
		if err := l.Result.Apply(ctx, rw); err != nil {
			return err
		}
	}

	l.logCallback()

	return nil
}

func (r *logger) Inject(flogger flamingo.Logger, responder *web.Responder) {
	r.logger = flogger
	r.responder = responder
}

// Filter a web request
func (r *logger) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	start := time.Now()

	webResponse := chain.Next(ctx, req, w)

	if webResponse == nil {
		webResponse = r.responder.ServerError(errors.New("nil response"))
	}

	return &loggedResponse{
		Result: webResponse,
		logCallback: func() {
			duration := time.Since(start)

			var cp func(msg interface{}, styles ...string) string
			switch {
			//case webResponse.GetStatus() >= 200 && webResponse.GetStatus() < 300:
			//	cp = color.Green
			//case webResponse.GetStatus() >= 300 && webResponse.GetStatus() < 400:
			//	cp = color.Blue
			//case webResponse.GetStatus() >= 400 && webResponse.GetStatus() < 500:
			//	cp = color.Yellow
			//case webResponse.GetStatus() >= 500 && webResponse.GetStatus() < 600:
			//	cp = color.Red
			default:
				cp = color.Grey
			}

			var extra string

			//if rr, ok := webResponse.(*web.RedirectResponse); ok {
			//	extra += " -> " + rr.Location
			//}

			//if r, ok := webResponse.(web.ErrorResponse); ok && r.Error != nil {
			//	extra += strings.Split(fmt.Sprintf(` | Error: %s`, r.Error.Error()), "\n")[0]
			//}

			var sizeStr string
			//if webResponse.GetContentLength() > 99999 {
			//	sizeStr = strconv.Itoa(webResponse.GetContentLength()/1000) + "kb"
			//} else {
			//	sizeStr = strconv.Itoa(webResponse.GetContentLength()) + "b"
			//}

			l := r.logger.
				WithContext(ctx).
				WithFields(
					map[flamingo.LogKey]interface{}{
						flamingo.LogKeyAccesslog: 1,
						//flamingo.LogKeyResponseCode: webResponse.GetStatus(),
						flamingo.LogKeyResponseTime: duration,
					},
				)
			defer l.Flush()

			l.Info(
				fmt.Sprintf(
					cp("%03d | %-8s | % 15s | % 7s | %s%s"),
					//webResponse.GetStatus(),
					//req.Request().Method,
					duration,
					sizeStr,
					//req.Request().RequestURI,
					extra,
				),
			)
		},
	}
}
