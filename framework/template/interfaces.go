package template

import (
	"io"

	"flamingo.me/flamingo/framework/web"
)

type (
	// Engine defines the basic template engine
	Engine interface {
		Render(context web.Context, name string, data interface{}) (io.Reader, error)
	}
)
