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

	// RenderPartials is used for progressive enhancements / rendering of partial template areas
	// usually this is requested via the appropriate javascript headers and taken care of in the framework renderer
	PartialEngine interface {
		RenderPartials(ctx context.Context, templateName string, data interface{}, partials []string) (map[string]string, error)
	}
)
