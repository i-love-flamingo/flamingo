package requestlogger

import (
	"time"

	"github.com/labstack/gommon/color"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	logger struct {
		Logger flamingo.Logger `inject:""`
	}
)

func (r *logger) Filter(ctx web.Context, chain *router.FilterChain) web.Response {
	start := time.Now()
	webResponse := chain.Next(ctx)

	duration := time.Since(start)
	webResponse.GetStatus()

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

	/*TODO error handling
	if event.Error != nil {
		extra += strings.Split(fmt.Sprintf(` | Error: %s`, event.Error), "\n")[0]
	}*/

	r.Logger.Printf(
		cp("%03d | %-8s | % 15s | % 6d byte | %s%s"),
		webResponse.GetStatus(),
		ctx.Request().Method,
		duration,
		0, //response.Size,
		ctx.Request().RequestURI,
		extra,
	)

	return webResponse
}
