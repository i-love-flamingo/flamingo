package interfaces

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo/v3/framework/web"
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
		{
			name: "render auth form on empty form submit",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
			},
			args: args{
				r: web.CreateRequest(
					&http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							RawQuery: "broker=testBroker",
							Scheme:   "http",
						},
						PostForm: url.Values{},
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
		{
			name: "renders auth form with message on invalid form data",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
			},
			args: args{
				r: web.CreateRequest(
					&http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							RawQuery: "broker=testBroker",
							Scheme:   "http",
						},
						PostForm: url.Values{
							"invalid-field": []string{"just to get into the form handling"},
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.Response{
				Status: http.StatusOK,
				Body: func() *bytes.Buffer {
					result := strings.Replace(defaultIDPTemplate, "{{.FormURL}}", "/test", 1)
					result = strings.Replace(result, "{{.Message}}", errMissingUsername, 1)

					return bytes.NewBuffer([]byte(result))
				}(),
				Header: http.Header{
					"ContentType": {"text/html; charset=utf-8"},
				},
			},
		},
		{
			name: "show error message on invalid user",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
			},
			args: args{
				r: web.CreateRequest(
					&http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							RawQuery: "broker=testBroker",
							Scheme:   "http",
						},
						PostForm: url.Values{
							"username": []string{"invalid user"},
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.Response{
				Status: http.StatusOK,
				Body: func() *bytes.Buffer {
					result := strings.Replace(defaultIDPTemplate, "{{.FormURL}}", "/test", 1)
					result = strings.Replace(result, "{{.Message}}", errInvalidUser, 1)

					return bytes.NewBuffer([]byte(result))
				}(),
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

			ctx := web.ContextWithSession(context.Background(), web.EmptySession())
			got := c.Auth(ctx, tt.args.r)

			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(bytes.Buffer{}, web.RouteRedirectResponse{})); diff != "" {
				t.Errorf("Auth() = %v, -got +want", diff)
			}
		})
	}
}
