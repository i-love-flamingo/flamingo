package flamingo

import (
	"context"
	"io"
)

type (
	// TemplateEngine defines the basic template engine
	TemplateEngine interface {
		Render(context context.Context, name string, data interface{}) (io.Reader, error)
	}

	// PartialTemplateEngine is used for progressive enhancements / rendering of partial template areas
	// usually this is requested via the appropriate javascript headers and taken care of in the framework renderer
	PartialTemplateEngine interface {
		RenderPartials(ctx context.Context, templateName string, data interface{}, partials []string) (map[string]io.Reader, error)
	}
)
