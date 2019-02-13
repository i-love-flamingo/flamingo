package domain

import "net/http"

type (
	// Handler responds to a request
	Handler interface {
		http.Handler
	}

	// HandlerProvider is a list of system HTTP handlers by path
	HandlerProvider func() map[string]Handler
)
