package interfaces

import (
	"net/http"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/google/go-cmp/cmp"
)

func TestIdentifier_Whitebox_Authenticate(t *testing.T) {
	type fields struct {
		responder     *web.Responder
		broker        string
		reverseRouter web.ReverseRouter
		eventRouter   flamingo.EventRouter
	}

	tests := []struct {
		name   string
		fields fields
		want   web.Result
	}{
		{
			name: "redirects to fake auth with broker parameter",
			fields: fields{
				responder:     &web.Responder{},
				broker:        "testBroker",
				reverseRouter: &mockRouter{},
				eventRouter:   &flamingo.DefaultEventRouter{},
			},
			want: &web.URLRedirectResponse{
				Response: web.Response{
					Status: http.StatusSeeOther,
					Header: http.Header{},
				},
				URL: &url.URL{
					Path: FakeAuthURL,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Identifier{
				responder:     tt.fields.responder,
				broker:        tt.fields.broker,
				reverseRouter: tt.fields.reverseRouter,
			}

			got := i.Authenticate(nil, nil)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Authenticate() -got +want: %v", diff)
			}
		})
	}
}
