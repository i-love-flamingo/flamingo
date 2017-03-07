package responder

import (
	"flamingo/core/flamingo/web"
	"flamingo/core/product/models"
)

type ErrorAware struct {
	DebugMode    string `inject:"param:debug.mode"`
	*RenderAware `inject:""`
}

type ErrorViewData struct {
	Error models.AppError
}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *ErrorAware) RenderError(context web.Context, error models.AppError) *web.ContentResponse {
	tpl := "pages/error"

	data := ErrorViewData{}
	data.Error = error

	if r.DebugMode == "0" {
		// Drop Message, should not be shown in public/prod env
		data.Error.Message = ""
	}

	response := r.RenderAware.Render(
		context,
		tpl,
		data,
	)

	response.Status = error.Code

	return response
}
