package router

import (
	"net/http"

	"go.aoe.com/flamingo/framework/web"
)

type (
	// OnRequestEvent contains the bare request
	OnRequestEvent struct {
		ResponseWriter http.ResponseWriter
		Request        *http.Request
		Ctx            web.Context
	}

	// OnResponseEvent is the event associated to OnResponse
	OnResponseEvent struct {
		Controller     Controller
		Response       web.Response
		Request        *http.Request
		ResponseWriter http.ResponseWriter
		Ctx            web.Context
	}

	// OnFinishEvent is the event object associated to OnFinish
	OnFinishEvent struct {
		ResponseWriter http.ResponseWriter
		Request        *http.Request
		Error          interface{}
		Ctx            web.Context
	}
)
