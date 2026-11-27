package responder_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Not found, forbidden and unavailable results render the template of their "flamingo.template.err*" key.
func TestRender_ErrorResultsUseTemplateKeys(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	// each key points to another template than its default
	swapped := []string{
		"flamingo.template.err404: error/403",
		"flamingo.template.err403: error/503",
		"flamingo.template.err503: error/404",
	}

	tests := []struct {
		name       string
		path       string
		config     []string
		wantStatus int
		wantBody   string
	}{
		{name: "render", path: "/page", wantStatus: http.StatusOK, wantBody: "page"},
		{name: "not found", path: "/notfound", wantStatus: http.StatusNotFound, wantBody: "error-404"},
		{name: "forbidden", path: "/forbidden", wantStatus: http.StatusForbidden, wantBody: "error-403"},
		{name: "unavailable", path: "/unavailable", wantStatus: http.StatusServiceUnavailable, wantBody: "error-503"},
		{name: "not found, set key", path: "/notfound", config: swapped, wantStatus: http.StatusNotFound, wantBody: "error-403"},
		{name: "forbidden, set key", path: "/forbidden", config: swapped, wantStatus: http.StatusForbidden, wantBody: "error-503"},
		{name: "unavailable, set key", path: "/unavailable", config: swapped, wantStatus: http.StatusServiceUnavailable, wantBody: "error-404"},
	}

	for _, tt := range tests { //nolint:paralleltest // boots an app, dingo scopes are process-global
		t.Run(tt.name, func(t *testing.T) {
			app := bootResponder(t, tt.config...)

			resp := app.Get(t, tt.path)

			assert.Equal(t, tt.wantStatus, resp.Status)
			assert.Contains(t, resp.Body, tt.wantBody)
		})
	}
}
