package fake

import (
	"bytes"
	"context"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type (
	mockRouter struct {
		mock.Mock
	}
)

var _ web.ReverseRouter = (*mockRouter)(nil)

func (m mockRouter) Relative(to string, params map[string]string) (*url.URL, error) {
	panic("implement me")
}

func (m mockRouter) Absolute(r *web.Request, to string, params map[string]string) (*url.URL, error) {
	resultURL := &url.URL{}

	result, _ := resultURL.Parse("test")

	return result, nil
}

func Test_idpController_Auth(t *testing.T) {
	type fields struct {
		responder     *web.Responder
		reverseRouter web.ReverseRouter
		template      string
	}

	type args struct {
		r *web.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   web.Result
	}{
		{
			name: "simple controller call with empty template to test fallback template",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
			},
			args: args{
				r: web.CreateRequest(
					&http.Request{
						Method: http.MethodGet,
						URL: &url.URL{
							RawQuery: "broker=testBroker",
							Scheme:   "http",
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.Response{
				Status: http.StatusOK,
				Body:   bytes.NewBuffer([]byte(strings.Replace(defaultIDPTemplate, "{{.FormURL}}", "/test", 1))),
				Header: http.Header{
					"ContentType": {"text/html; charset=utf-8"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &idpController{
				responder:     tt.fields.responder,
				reverseRouter: tt.fields.reverseRouter,
				template:      tt.fields.template,
			}

			got := c.Auth(context.Background(), tt.args.r)

			if diff := cmp.Diff(got, tt.want, cmp.Options{cmpopts.IgnoreUnexported(bytes.Buffer{})}); diff != "" {
				t.Errorf("Auth() = %v, -got +want", diff)
			}
		})
	}
}
