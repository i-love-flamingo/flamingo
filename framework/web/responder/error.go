package responder

import (
	"flamingo/framework/web"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// ErrorAware responder trait
	ErrorAware interface {
		Error(context web.Context, err error) web.Response
		ErrorNotFound(context web.Context, error error) web.Response
	}

	// FlamingoErrorAware responder can return errors
	FlamingoErrorAware struct {
		DebugMode   bool `inject:"config:debug.mode"`
		RenderAware `inject:""`
	}

	// ErrorViewData for template rendering
	ErrorViewData struct {
		Code  int
		Error error
	}

	// DebugError holds additional information
	DebugError struct {
		Err error
	}

	// EmptyError in case we want to hide our error
	EmptyError struct{}

	st interface {
		StackTrace() errors.StackTrace
	}
)

var _ ErrorAware = &FlamingoErrorAware{}

// Error implements error interface
func (de DebugError) Error() string {
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

	response.(*web.ContentResponse).Status = http.StatusNotFound

	return response
}

// Error returns a web.ContentResponse with status 503 and ContentType text/html
func (r *FlamingoErrorAware) Error(context web.Context, err error) web.Response {
	var response web.Response

	if !r.DebugMode {
		response = r.RenderAware.Render(
			context,
			"pages/error/503",
			ErrorViewData{Error: EmptyError{}, Code: 500},
		)
	} else {
		response = r.RenderAware.Render(
			context,
			"pages/error/503",
			ErrorViewData{Error: DebugError{err}, Code: 500},
		)
	}

	response.(*web.ContentResponse).Status = http.StatusInternalServerError

	return response
}
