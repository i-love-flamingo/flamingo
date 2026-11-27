package responder_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirect_RouteAndURL(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	tests := []struct {
		path         string
		wantLocation string
	}{
		{path: "/redirect/route", wantLocation: "/product/1"},
		{path: "/redirect/url", wantLocation: "https://example.com/target"},
	}

	app := bootResponder(t)

	for _, tt := range tests { //nolint:paralleltest // boots an app, dingo scopes are process-global
		t.Run(tt.path, func(t *testing.T) {
			resp := app.Get(t, tt.path)

			assert.Equal(t, http.StatusSeeOther, resp.Status)
			assert.Equal(t, tt.wantLocation, resp.Header.Get("Location"))
		})
	}
}
