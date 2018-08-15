package responder

import (
	"context"
	"fmt"
	"net/http"

	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (
	// ErrorAware responder trait
	ErrorAware interface {
		Error(context context.Context, err error) web.Response
		ErrorNotFound(context context.Context, err error) web.Response
		ErrorWithCode(context context.Context, err error, httpStatus int) web.Response
	}

	// FlamingoErrorAware responder can return errors
	FlamingoErrorAware struct {
		RenderAware      `inject:""`
		DebugMode        bool            `inject:"config:debug.mode"`
		Tpl404           string          `inject:"config:flamingo.template.err404"`
		TplErrorWithCode string          `inject:"config:flamingo.template.errWithCode"`
		Tpl503           string          `inject:"config:flamingo.template.err503"`
		Logger           flamingo.Logger `inject:""`
	}

	// ErrorViewData for template rendering
	ErrorViewData struct {
		Code      int    `json:"code"`
		Error     error  `json:"error"`
		Errortext string `json:"errortex"`
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
func (r *FlamingoErrorAware) ErrorNotFound(context context.Context, error error) web.Response {
	var response web.Response
	r.Logger.WithField("category", "error_aware").WithField("errorTemplate", r.Tpl404).Info("ErrorNotFound", error)
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

// ErrorWithCode returns a web.ContentResponse with give status code and ContentType text/html
func (r *FlamingoErrorAware) ErrorWithCode(context context.Context, error error, httpStatus int) web.Response {
	var response web.Response
	r.Logger.WithField("category", "error_aware").WithField("errorTemplate", r.TplErrorWithCode).Error("Error with code ", httpStatus, error)
	if !r.DebugMode {
		response = r.RenderAware.Render(
			context,
			r.TplErrorWithCode,
			ErrorViewData{Error: EmptyError{}, Code: httpStatus, Errortext: http.StatusText(httpStatus)},
		)
	} else {
		response = r.RenderAware.Render(
			context,
			r.TplErrorWithCode,
			ErrorViewData{Error: DebugError{error}, Code: httpStatus, Errortext: http.StatusText(httpStatus)},
		)
	}

	response.(*web.ContentResponse).Status = httpStatus

	return web.ErrorResponse{Response: response, Error: error}
}

// Error returns a web.ContentResponse with status 503 and ContentType text/html
func (r *FlamingoErrorAware) Error(context context.Context, err error) web.Response {
	var response web.Response
	r.Logger.WithField("category", "error_aware").WithField("errorTemplate", r.Tpl503).Error("Error ", err.Error())
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
			ErrorViewData{Error: DebugError{err}, Code: 500},
		)
	}

	if r, ok := response.(*web.ContentResponse); ok {
		r.Status = http.StatusInternalServerError
	}

	return web.ErrorResponse{Response: response, Error: err}
}
