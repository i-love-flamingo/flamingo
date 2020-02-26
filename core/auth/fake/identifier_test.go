package fake

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

type (
	mockRouter struct {
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
	return url.Parse(strings.ReplaceAll("/core/auth/login/:broker", ":broker", m.broker))
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
				r: buildRequest(http.MethodGet, nil),
			},
			want: wantFormResponseWithMessage(""),
		},
		{
			name: "render auth form on empty form submit",
			fields: fields{
				config: fakeConfig{},
			},
			args: args{
				r: buildRequest(http.MethodPost, &url.Values{}),
			},
			want: wantFormResponseWithMessage(""),
		},
		{
			name: "renders auth form with message on invalid form data",
			fields: fields{
				config: fakeConfig{},
			},
			args: args{
				r: buildRequest(http.MethodPost, &url.Values{
					"invalid-field": []string{"just to get into the form handling"},
				}),
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
				r: buildRequest(http.MethodPost, &url.Values{
					"username": []string{"nonexistent user"},
				}),
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
				r: buildRequest(http.MethodPost, &url.Values{
					"username": []string{"user_b"},
				}),
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
				r: buildRequest(http.MethodPost, &url.Values{
					"username": []string{"user_b"},
					"password": []string{"invalid password"},
				}),
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
				r: buildRequest(http.MethodPost, &url.Values{
					"username": []string{"user_b"},
					"password": []string{"test_b"},
				}),
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
			i := new(identifier)

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

func buildRequest(method string, formValues *url.Values) *web.Request {
	newRequest := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme: "http",
		},
	}

	if formValues != nil && http.MethodPost == method {
		newRequest.PostForm = *formValues
	}

	result := web.CreateRequest(
		newRequest,
		web.EmptySession(),
	)

	result.Params = web.RequestParams{"broker": "testBroker"}

	return result
}

func wantFormResponseWithMessage(message string) *web.Response {

	t := template.New("fake")
	t, err := t.Parse(defaultLoginTemplate)
	if err != nil {
		return i.responder.ServerError(err)
	}

	var body = new(bytes.Buffer)
	var errMsg string

	if formError != nil {
		errMsg = formError.Error()
	}

	err = t.Execute(
		body,
		viewData{
			FormURL:    formURL.String(),
			Message:    errMsg,
			UsernameID: i.config.UsernameFieldID,
			PasswordID: i.config.PasswordFieldID,
		})
	if err != nil {
		return i.responder.ServerError(err)
	}

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
