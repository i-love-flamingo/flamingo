package responder

import (
	"flamingo/framework/web"
	"fmt"
	"net/http"
)

type (
	ErrorAware struct {
		DebugMode    bool `inject:"config:debug.mode"`
		*RenderAware `inject:""`
	}

	ErrorViewData struct {
		Code  int
		Error error
	}

	DebugError struct {
		Err error
	}

	EmptyError struct{}
)

func (de DebugError) Error() string {
	return fmt.Sprintf("%+v", de.Err)
}

func (ee EmptyError) Error() string {
	return ""
}

// ErrorNotFound returns a web.ContentResponse with status 503 and ContentType text/html
func (r *ErrorAware) ErrorNotFound(context web.Context, error error) *web.ContentResponse {
	var response *web.ContentResponse

	if !r.DebugMode {
		response = r.RenderAware.Render(
			context,
			"pages/error/404",
			ErrorViewData{Error: EmptyError{}, Code: 404},
		)
	} else {
		response = r.RenderAware.Render(
			context,
			"pages/error/404",
			ErrorViewData{Error: DebugError{error}, Code: 404},
		)
	}

	response.Status = http.StatusNotFound

	return response
}

// Error returns a web.ContentResponse with status 503 and ContentType text/html
func (r *ErrorAware) Error(context web.Context, error error) *web.ContentResponse {
	var response *web.ContentResponse

	if !r.DebugMode {
		response = r.RenderAware.Render(
			context,
			"pages/error/503",
			ErrorViewData{Error: EmptyError{}, Code: 503},
		)
	} else {
		response = r.RenderAware.Render(
			context,
			"pages/error/503",
			ErrorViewData{Error: DebugError{error}, Code: 503},
		)
	}

	response.Status = http.StatusServiceUnavailable

	return response
}
