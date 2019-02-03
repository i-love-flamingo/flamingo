package domain

import (
	"net/http"
)

type (
	// Service to retrieve the base URL
	Service interface {
		BaseURL() string
		BaseDomain() string
		DetermineBase(r *http.Request) string
	}
)
