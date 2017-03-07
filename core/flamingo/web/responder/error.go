package responder

import (
	"flamingo/core/flamingo/web"
	"flamingo/core/product/models"
)

type ErrorAware struct {
	*RenderAware `inject:""`
}

type ErrorAwareDebug struct {
	*RenderAware `inject:""`
}

type ErrorViewData struct {
	Error models.AppError
}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *ErrorAware) RenderError(context web.Context, error models.AppError) *web.ContentResponse {
	tpl := "pages/error"

	data := ErrorViewData{}

	// Drop Message, should not be shown in public/prod env
	data.Error.Message = ""

	data.Error = error

	response := r.RenderAware.Render(
		context,
		tpl,
		data,
	)

	return response
}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *ErrorAwareDebug) RenderError(context web.Context, error models.AppError) *web.ContentResponse {
	tpl := "pages/error"

	data := ErrorViewData{}

	data.Error = error

	response := r.RenderAware.Render(
		context,
		tpl,
		data,
	)

	return response
}
