package routing_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/catalog"
)

func TestRoutes_ParamsAndDefaults(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootAreas(t, nil, sameModules(new(catalog.Module)), areaFiles())

	tests := []struct {
		path string
		want string
	}{
		{path: "/en/product/7", want: `{"id":"7"}`},
		{path: "/en/list", want: `{"sort":"name"}`},
		{path: "/en/list?sort=price", want: `{"sort":"price"}`},
	}

	for _, tt := range tests { //nolint:paralleltest // boots an app, dingo scopes are process-global
		t.Run(tt.path, func(t *testing.T) {
			resp := app.Get(t, tt.path)

			assert.Equal(t, http.StatusOK, resp.Status)
			assert.JSONEq(t, tt.want, resp.Body)
		})
	}
}
