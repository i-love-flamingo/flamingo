package web

import (
	"net/http"
)

type (
	// OnRequestEvent contains the bare request
	OnRequestEvent struct {
		Request        *Request
		ResponseWriter http.ResponseWriter
	}

	// OnResponseEvent is the event associated to OnResponse
	OnResponseEvent struct {
		OnRequestEvent
		Result Result
	}

	// OnFinishEvent is the event object associated to OnFinish
	OnFinishEvent struct {
		OnRequestEvent
		Error error
	}
)
