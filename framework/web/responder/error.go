package responder

import (
	"flamingo/framework/web"
)

type ErrorAware struct {
	DebugMode    bool `inject:"config:debug.mode"`
	*RenderAware `inject:""`
}

type ErrorViewData struct {
	//Error models.AppError
}

// RenderError returns a web.ContentResponse with status 500 and ContentType text/html
func (r *ErrorAware) RenderError(context web.Context, error interface{}) *web.ContentResponse {
	tpl := "pages/error"

	data := ErrorViewData{}
	//data.Error = error

	if !r.DebugMode {
		// Drop Message, should not be shown in public/prod env
		//data.Error.Message = ""
	}

	response := r.RenderAware.Render(
		context,
		tpl,
		data,
	)

	/*
		Response.Status only modified if error.Code is a valid HTTP Status Code
		otherwise 200 is kept from the default response as its probably a proprietary
		error code (which is also ok, just for display)
	*/
	/*
		if error.Code > 99 && error.Code < 1000 {
			response.Status = error.Code
		}
	*/

	return response
}
