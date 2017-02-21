package template

import (
	"flamingo/core/flamingo/web"
	"io"
)

type (
	Engine interface {
		Render(context web.Context, name string, data interface{}) io.Reader
	}
)
