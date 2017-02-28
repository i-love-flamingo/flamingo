package router

import (
	"flamingo/core/flamingo/web"
	"net/http"
)

type (
	// Event is the type for router event keys
	Event string

	// OnRequest is fired whenever a request is handled
	OnRequest interface {
		OnRequest(event *OnRequestEvent)
	}

	// OnRequestEvent contains the bare request
	OnRequestEvent struct {
		ResponseWriter http.ResponseWriter
		Request        *http.Request
	}

	// OnResponse is fired whenever a request was handled, right before the apply()
	// Note: this does not fire for http.Handler!
	OnResponse interface {
		OnResponse(event *OnResponseEvent)
	}

	// OnResponseEvent is the event associated to OnResponse
	OnResponseEvent struct {
		Controller Controller
		Response   web.Response
	}

	// OnFinish is fired when the request finished, usually the response has already been written
	OnFinish interface {
		OnFinish(event *OnFinishEvent)
	}

	// OnFinishEvent is the event object associated to OnFinish
	OnFinishEvent struct {
		ResponseWriter http.ResponseWriter
		Request        *http.Request
		Error          interface{}
	}
)

const (
	REQUEST  Event = "router.request"
	RESPONSE       = "router.response"
	FINISH         = "router.finish"
)

// Dispatch the OnRequest event
func (event *OnRequestEvent) Dispatch(subscriber interface{}) {
	subscriber.(OnRequest).OnRequest(event)
}

// Dispatch the OnResponse event
func (event *OnResponseEvent) Dispatch(subscriber interface{}) {
	subscriber.(OnResponse).OnResponse(event)
}

// Dispatch the OnFinish event
func (event *OnFinishEvent) Dispatch(subscriber interface{}) {
	subscriber.(OnFinish).OnFinish(event)
}
