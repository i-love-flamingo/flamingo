package template

import (
	"flamingo/core/flamingo/web"
	"io"
)

type (
	// Engine defines the basic template engine
	Engine interface {
		Render(context web.Context, name string, data interface{}) io.Reader
	}
)
