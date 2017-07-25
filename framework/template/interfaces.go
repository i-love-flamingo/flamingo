package template

import (
	"flamingo/framework/web"
	"io"
)

type (
	// Engine defines the basic template engine
	Engine interface {
		Render(context web.Context, name string, data interface{}) (io.Reader, error)
	}
)
