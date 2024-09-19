package http

import (
	"net/http"
)

type (
	HandlerWrapper func(http.Handler) http.Handler
)
