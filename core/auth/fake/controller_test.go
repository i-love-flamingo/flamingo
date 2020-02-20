package fake

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"flamingo.me/flamingo/v3/framework/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	mockRouter struct {
		mock.Mock
		broker string
	}
)

var _ web.ReverseRouter = (*mockRouter)(nil)

// Relative mock action
func (m *mockRouter) Relative(to string, params map[string]string) (*url.URL, error) {
	panic("not implemented")
}

// Absolut mock action
func (m *mockRouter) Absolute(r *web.Request, to string, params map[string]string) (*url.URL, error) {
	resultURL := &url.URL{}

	return resultURL.Parse(strings.ReplaceAll(FakeAuthURL, ":broker", m.broker))
}

func Test_idpController_Auth(t *testing.T) {
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
			name: "simple controller call with empty template to test fallback template",
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
						"user_a": userConfig{
							Password: "test_a",
							Otp:      "otp_a",
						},
						"user_b": userConfig{
							Password: "test_b",
							Otp:      "otp_b",
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
			name: "login for valid user without pwd/otp",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": userConfig{
							Password: "test_a",
							Otp:      "otp_a",
						},
						"user_b": userConfig{
							Password: "test_b",
							Otp:      "otp_b",
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
			name: "login for valid user with password mismatch but w/o otp",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": userConfig{
							Password: "test_a",
							Otp:      "otp_a",
						},
						"user_b": userConfig{
							Password: "test_b",
							Otp:      "otp_b",
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
			name: "login for valid user / valid password / w/o otp",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": userConfig{
							Password: "test_a",
							Otp:      "otp_a",
						},
						"user_b": userConfig{
							Password: "test_b",
							Otp:      "otp_b",
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
		{
			name: "login for valid user with password, otp mismatch",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": userConfig{
							Password: "test_a",
							Otp:      "otp_a",
						},
						"user_b": userConfig{
							Password: "test_b",
							Otp:      "otp_b",
						},
					},
					ValidatePassword: true,
					ValidateOtp:      true,
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
								"otp":      []string{"invalid otp"},
							},
						},
						web.EmptySession(),
					),
				),
			},
			want: wantFormResponseWithMessage(errOtpMismatch.Error()),
		},
		{
			name: "login for valid user / valid password / valid otp",
			fields: fields{
				config: fakeConfig{
					UserConfig: map[string]userConfig{
						"user_a": userConfig{
							Password: "test_a",
							Otp:      "otp_a",
						},
						"user_b": userConfig{
							Password: "test_b",
							Otp:      "otp_b",
						},
					},
					ValidatePassword: true,
					ValidateOtp:      true,
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
								"otp":      []string{"otp_b"},
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
			c := new(controller)

			c.Inject(&web.Responder{}, &mockRouter{broker: "testBroker"})

			// prepare identifier with default config values not required by the test
			tt.fields.config.Broker = "testBroker"
			tt.fields.config.UsernameFieldID = "username"
			tt.fields.config.PasswordFieldID = "password"
			tt.fields.config.OtpFieldID = "otp"
			identifierConfig[tt.fields.config.Broker] = tt.fields.config

			got := c.Auth(web.ContextWithSession(context.Background(), web.EmptySession()), tt.args.r)

			assert.Equal(t, tt.want, got)
		})
	}
}

func wantFormResponseWithMessage(message string) *web.Response {
	result := strings.ReplaceAll(defaultIDPTemplate, "{{.FormURL}}", "/core/auth/fake/testBroker")
	result = strings.ReplaceAll(result, "{{.Message}}", message)
	result = strings.ReplaceAll(result, "{{.UsernameID}}", defaultUserNameFieldID)
	result = strings.ReplaceAll(result, "{{.PasswordID}}", defaultPasswordFieldID)
	result = strings.ReplaceAll(result, "{{.OtpID}}", defaultOtpFieldID)

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
