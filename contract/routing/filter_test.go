package routing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/catalog"
	"flamingo.me/flamingo/v3/framework/web"
)

// orderFilter adds its name to the header "X-Filter-Order".
type orderFilter string

func (f orderFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	w.Header().Add("X-Filter-Order", string(f))

	return chain.Next(ctx, req, w)
}

// Two module types, because an area loads each module type once.
type (
	firstFilterModule  struct{}
	secondFilterModule struct{}
)

func (*firstFilterModule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).ToInstance(orderFilter("first"))
}

func (*secondFilterModule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).ToInstance(orderFilter("second"))
}

func TestFilter_RunsInModuleOrder(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	tests := []struct {
		name    string
		modules []dingo.Module
		want    []string
	}{
		{
			name:    "first, second",
			modules: []dingo.Module{new(firstFilterModule), new(secondFilterModule)},
			want:    []string{"first", "second"},
		},
		{
			name:    "second, first",
			modules: []dingo.Module{new(secondFilterModule), new(firstFilterModule)},
			want:    []string{"second", "first"},
		},
	}

	for _, tt := range tests { //nolint:paralleltest // boots an app, dingo scopes are process-global
		t.Run(tt.name, func(t *testing.T) {
			app := bootAreas(t, tt.modules, sameModules(new(catalog.Module)), areaFiles())

			resp := app.Get(t, "/en/product/1")

			assert.Equal(t, http.StatusOK, resp.Status)
			assert.Equal(t, tt.want, resp.Header.Values("X-Filter-Order"))
		})
	}
}
