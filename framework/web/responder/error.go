package responder

import (
	"fmt"
	"net/http"

	"go.aoe.com/flamingo/framework/web"
)

type (
	// ErrorAware responder trait
	ErrorAware interface {
		Error(context web.Context, err error) web.Response
		ErrorNotFound(context web.Context, err error) web.Response
		ErrorGone(context web.Context, error error) web.Response
	}

	// FlamingoErrorAware responder can return errors
	FlamingoErrorAware struct {
		RenderAware `inject:""`
		DebugMode   bool   `inject:"config:debug.mode"`
		Tpl404      string `inject:"config:flamingo.template.err404"`
		Tpl410      string `inject:"config:flamingo.template.err410"`
		Tpl503      string `inject:"config:flamingo.template.err503"`
	}

	// ErrorViewData for template rendering
	ErrorViewData struct {
		Code  int   `json:"code"`
		Error error `json:"error"`
	}

	// DebugError holds additional information
	DebugError struct {
		Err error
	}

	// EmptyError in case we want to hide our error
	EmptyError struct{}
)

var _ ErrorAware = &FlamingoErrorAware{}

// Error implements error interface
func (de DebugError) Error() string {
	return fmt.Sprintf("%+v", de.Err)
}

// Error implements error interface
func (de DebugError) String() string {
	return fmt.Sprintf("%+v", de.Err)
}

// Error implements error interface
func (ee EmptyError) Error() string {
	return ""
}

// ErrorNotFound returns a web.ContentResponse with status 404 and ContentType text/html
func (r *FlamingoErrorAware) ErrorNotFound(context web.Context, error error) web.Response {
	var response web.Response

	if !r.DebugMode {
		response = r.RenderAware.Render(
			context,
			r.Tpl404,
			ErrorViewData{Error: EmptyError{}, Code: 404},
		)
	} else {
		response = r.RenderAware.Render(
			context,
			r.Tpl404,
			ErrorViewData{Error: DebugError{error}, Code: 404},
		)
	}

	response.(*web.ContentResponse).Status = http.StatusNotFound

	return web.ErrorResponse{Response: response, Error: error}
}

// ErrorGone returns a web.ContentResponse with status 410 and ContentType text/html
func (r *FlamingoErrorAware) ErrorGone(context web.Context, error error) web.Response {
	var response web.Response

	if !r.DebugMode {
		response = r.RenderAware.Render(
			context,
			r.Tpl410,
			ErrorViewData{Error: EmptyError{}, Code: 410},
		)
	} else {
		response = r.RenderAware.Render(
			context,
			r.Tpl410,
			ErrorViewData{Error: DebugError{error}, Code: 410},
		)
	}

	response.(*web.ContentResponse).Status = http.StatusGone

	return web.ErrorResponse{Response: response, Error: error}
}

// Error returns a web.ContentResponse with status 503 and ContentType text/html
func (r *FlamingoErrorAware) Error(context web.Context, error error) web.Response {
	var response web.Response

	if !r.DebugMode {
		response = r.RenderAware.Render(
			context,
			r.Tpl503,
			ErrorViewData{Error: EmptyError{}, Code: 500},
		)
	} else {
		response = r.RenderAware.Render(
			context,
			r.Tpl503,
			ErrorViewData{Error: DebugError{error}, Code: 500},
		)
	}

	if r, ok := response.(*web.ContentResponse); ok {
		r.Status = http.StatusInternalServerError
	}

	return web.ErrorResponse{Response: response, Error: error}
}
