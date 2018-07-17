package template

import (
	"context"
	"io"
)

type (
	// Engine defines the basic template engine
	Engine interface {
		Render(context context.Context, name string, data interface{}) (io.Reader, error)
	}
)
