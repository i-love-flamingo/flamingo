package fake

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

type (
	mockRouter struct {
		mock.Mock
		broker string
	}
)

var _ web.ReverseRouter = (*mockRouter)(nil)

// Relative mock action
func (m *mockRouter) Relative(_ string, _ map[string]string) (*url.URL, error) {
	panic("not implemented")
}

// Absolute mock action
func (m *mockRouter) Absolute(_ *web.Request, _ string, _ map[string]string) (*url.URL, error) {
	resultURL := &url.URL{}

	return resultURL.Parse(strings.ReplaceAll("/core/auth/login/:broker", ":broker", m.broker))
}

func Test_Identifier_Authenticate(t *testing.T) {
	t.Parallel()

	type fields struct {
		config fakeConfig
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
			name: "simple identifier call with empty template to test fallback template",
			fields: fields{
				config: fakeConfig{},
			},
			args: args{
				r: addRequestParameters(web.CreateRequest(
					&http.Request{
						Method: http.MethodGet,
						URL: &url.URL{
							Scheme: "http",
						},
					},
					web.EmptySession(),
				),
				),
			},
			want: wantFormResponseWithMessage(""),
		},
		{
			name: "render auth form on empty form submit",
			fields: fields{
				config: fakeConfig{},
			},
			args: args{
				r: addRequestParameters(
					web.CreateRequest(
						&http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Scheme: "http",
							},
							PostForm: url.Values{},
						},
						web.EmptySession(),
					),
				),
			},
			want: wantFormResponseWithMessage(""),
		},
		{
			name: "renders auth form with message on invalid form data",
			fields: fields{
				config: fakeConfig{},
			},
			args: args{
				r: addRequestParameters(
					web.CreateRequest(
						&http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Scheme: "http",
							},
							PostForm: url.Values{
								"invalid-field": []string{"just to get into the form handling"},
							},
						},
						web.EmptySession(),
					),
				),
			},
			want: wantFormResponseWithMessage(errMissingUsername.Error()),
		},
		{
			name: "show error message on invalid user",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": {
							Password: "test_a",
						},
						"user_b": {
							Password: "test_b",
						},
					},
				},
			},
			args: args{
				r: addRequestParameters(
					web.CreateRequest(
						&http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Scheme: "http",
							},
							PostForm: url.Values{
								"username": []string{"nonexistent user"},
							},
						},
						web.EmptySession(),
					),
				),
			},
			want: wantFormResponseWithMessage(errInvalidUser.Error()),
		},
		{
			name: "login for valid user without pwd",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": {
							Password: "test_a",
						},
						"user_b": {
							Password: "test_b",
						},
					},
				},
			},
			args: args{
				r: addRequestParameters(
					web.CreateRequest(
						&http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Scheme: "http",
							},
							PostForm: url.Values{
								"username": []string{"user_b"},
							},
						},
						web.EmptySession(),
					),
				),
			},
			want: &web.RouteRedirectResponse{
				Response: web.Response{
					Status: http.StatusSeeOther,
					Header: http.Header{},
				},
				To:   "core.auth.callback",
				Data: map[string]string{"broker": "testBroker"},
			},
		},
		{
			name: "login for valid user with password mismatch",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": {
							Password: "test_a",
						},
						"user_b": {
							Password: "test_b",
						},
					},
					ValidatePassword: true,
				},
			},
			args: args{
				r: addRequestParameters(
					web.CreateRequest(
						&http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Scheme: "http",
							},
							PostForm: url.Values{
								"username": []string{"user_b"},
								"password": []string{"invalid password"},
							},
						},
						web.EmptySession(),
					),
				),
			},
			want: wantFormResponseWithMessage(errPasswordMismatch.Error()),
		},
		{
			name: "login for valid user / valid password",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": {
							Password: "test_a",
						},
						"user_b": {
							Password: "test_b",
						},
					},
					ValidatePassword: true,
				},
			},
			args: args{
				r: addRequestParameters(
					web.CreateRequest(
						&http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Scheme: "http",
							},
							PostForm: url.Values{
								"username": []string{"user_b"},
								"password": []string{"test_b"},
							},
						},
						web.EmptySession(),
					),
				),
			},
			want: &web.RouteRedirectResponse{
				Response: web.Response{
					Status: http.StatusSeeOther,
					Header: http.Header{},
				},
				To:   "core.auth.callback",
				Data: map[string]string{"broker": "testBroker"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := new(Identifier)

			i.Inject(&mockRouter{broker: "testBroker"}, &web.Responder{}, &flamingo.DefaultEventRouter{})

			// prepare identifier with default config values not required by the test

			i.config = tt.fields.config
			i.broker = "testBroker"
			i.config.UsernameFieldID = "username"
			i.config.PasswordFieldID = "password"

			got := i.Authenticate(web.ContextWithSession(context.Background(), web.EmptySession()), tt.args.r)

			assert.Equal(t, tt.want, got)
		})
	}
}

func wantFormResponseWithMessage(message string) *web.Response {
	result := strings.ReplaceAll(defaultLoginTemplate, "{{.FormURL}}", "/core/auth/login/testBroker")
	result = strings.ReplaceAll(result, "{{.Message}}", message)
	result = strings.ReplaceAll(result, "{{.UsernameID}}", defaultUserNameFieldID)
	result = strings.ReplaceAll(result, "{{.PasswordID}}", defaultPasswordFieldID)

	return &web.Response{
		Status: http.StatusOK,
		Body:   bytes.NewBufferString(result),
		Header: http.Header{
			"ContentType": {"text/html; charset=utf-8"},
		},
	}
}

func addRequestParameters(request *web.Request) *web.Request {
	request.Params = web.RequestParams{"broker": "testBroker"}

	return request
}
