package interfaces

import (
	"bytes"
	"context"
	"errors"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	mockRouter struct {
		mock.Mock
	}
)

var _ web.ReverseRouter = (*mockRouter)(nil)

// Relative mock action
func (m *mockRouter) Relative(to string, params map[string]string) (*url.URL, error) {
	_ = to
	_ = params

	return nil, errors.New("not implemented")
}

// Absolut mock action
func (m *mockRouter) Absolute(r *web.Request, to string, params map[string]string) (*url.URL, error) {
	_ = r
	_ = to
	_ = params

	resultURL := &url.URL{}
	result, _ := resultURL.Parse("test")

	return result, nil
}

func Test_idpController_Auth(t *testing.T) {
	t.Parallel()

	type fields struct {
		responder        *web.Responder
		reverseRouter    web.ReverseRouter
		template         string
		fakeUserData     config.Map
		validatePassword bool
		validateOtp      bool
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
				Body: func() *bytes.Buffer {
					result := strings.Replace(defaultIDPTemplate, "{{.FormURL}}", "/test", 1)
					result = strings.Replace(result, "{{.Message}}", "", 1)
					result = replaceStandardFormIDs(result)

					return bytes.NewBuffer([]byte(result))
				}(),
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
				Body: func() *bytes.Buffer {
					result := strings.Replace(defaultIDPTemplate, "{{.FormURL}}", "/test", 1)
					result = strings.Replace(result, "{{.Message}}", "", 1)
					result = replaceStandardFormIDs(result)

					return bytes.NewBuffer([]byte(result))
				}(),
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
					result = replaceStandardFormIDs(result)

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
				fakeUserData: config.Map{
					"user_a": config.Map{
						"password": "test_a",
						"otp":      "otp_a",
					},
					"user_b": config.Map{
						"password": "test_b",
						"otp":      "otp_b",
					},
				},
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
							"username": []string{"nonexistent user"},
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
					result = replaceStandardFormIDs(result)

					return bytes.NewBuffer([]byte(result))
				}(),
				Header: http.Header{
					"ContentType": {"text/html; charset=utf-8"},
				},
			},
		},
		{
			name: "login for valid user without pwd/otp",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
				fakeUserData: config.Map{
					"user_a": config.Map{
						"password": "test_a",
						"otp":      "otp_a",
					},
					"user_b": config.Map{
						"password": "test_b",
						"otp":      "otp_b",
					},
				},
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
							"username": []string{"user_b"},
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.RouteRedirectResponse{
				Response: web.Response{
					Status: http.StatusSeeOther,
					Header: http.Header{},
				},
				To:   "core.auth.callback(broker)",
				Data: map[string]string{"broker": "testBroker"},
			},
		},
		{
			name: "login for valid user with password mismatch but w/o otp",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
				fakeUserData: config.Map{
					"user_a": config.Map{
						"password": "test_a",
						"otp":      "otp_a",
					},
					"user_b": config.Map{
						"password": "test_b",
						"otp":      "otp_b",
					},
				},
				validatePassword: true,
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
							"username": []string{"user_b"},
							"password": []string{"invalid password"},
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.Response{
				Status: http.StatusOK,
				Body: func() *bytes.Buffer {
					result := strings.Replace(defaultIDPTemplate, "{{.FormURL}}", "/test", 1)
					result = strings.Replace(result, "{{.Message}}", errPasswordMismatch, 1)
					result = replaceStandardFormIDs(result)

					return bytes.NewBuffer([]byte(result))
				}(),
				Header: http.Header{
					"ContentType": {"text/html; charset=utf-8"},
				},
			},
		},
		{
			name: "login for valid user / valid password / w/o otp",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
				fakeUserData: config.Map{
					"user_a": config.Map{
						"password": "test_a",
						"otp":      "otp_a",
					},
					"user_b": config.Map{
						"password": "test_b",
						"otp":      "otp_b",
					},
				},
				validatePassword: true,
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
							"username": []string{"user_b"},
							"password": []string{"test_b"},
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.RouteRedirectResponse{
				Response: web.Response{
					Status: http.StatusSeeOther,
					Header: http.Header{},
				},
				To:   "core.auth.callback(broker)",
				Data: map[string]string{"broker": "testBroker"},
			},
		},
		{
			name: "login for valid user with password, otp mismatch",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
				fakeUserData: config.Map{
					"user_a": config.Map{
						"password": "test_a",
						"otp":      "otp_a",
					},
					"user_b": config.Map{
						"password": "test_b",
						"otp":      "otp_b",
					},
				},
				validatePassword: true,
				validateOtp:      true,
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
							"username": []string{"user_b"},
							"password": []string{"test_b"},
							"m2fa-otp": []string{"invalid otp"},
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.Response{
				Status: http.StatusOK,
				Body: func() *bytes.Buffer {
					result := strings.Replace(defaultIDPTemplate, "{{.FormURL}}", "/test", 1)
					result = strings.Replace(result, "{{.Message}}", errOtpMismatch, 1)
					result = replaceStandardFormIDs(result)

					return bytes.NewBuffer([]byte(result))
				}(),
				Header: http.Header{
					"ContentType": {"text/html; charset=utf-8"},
				},
			},
		},
		{
			name: "login for valid user / valid password / valid otp",
			fields: fields{
				responder:     &web.Responder{},
				reverseRouter: &mockRouter{},
				template:      "",
				fakeUserData: config.Map{
					"user_a": config.Map{
						"password": "test_a",
						"otp":      "otp_a",
					},
					"user_b": config.Map{
						"password": "test_b",
						"otp":      "otp_b",
					},
				},
				validatePassword: true,
				validateOtp:      true,
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
							"username": []string{"user_b"},
							"password": []string{"test_b"},
							"m2fa-otp": []string{"otp_b"},
						},
					},
					web.EmptySession(),
				),
			},
			want: &web.RouteRedirectResponse{
				Response: web.Response{
					Status: http.StatusSeeOther,
					Header: http.Header{},
				},
				To:   "core.auth.callback(broker)",
				Data: map[string]string{"broker": "testBroker"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(IdpController)

			testConfig := new(struct {
				Template         string
				UserConfig       config.Map
				ValidatePassword bool
				ValidateOtp      bool
				UsernameFieldID  string
				PasswordFieldID  string
				OtpFieldID       string
			})

			testConfig.Template = tt.fields.template
			testConfig.ValidatePassword = tt.fields.validatePassword
			testConfig.ValidateOtp = tt.fields.validateOtp

			if tt.fields.fakeUserData != nil {
				testConfig.UserConfig = tt.fields.fakeUserData
			}

			c.Inject(
				tt.fields.responder,
				tt.fields.reverseRouter,
				(*struct {
					Template         string     `inject:"config:core.auth.fake.loginTemplate,optional"`
					UserConfig       config.Map `inject:"config:core.auth.fake.userConfig"`
					ValidatePassword bool       `inject:"config:core.auth.fake.validatePassword,optional"`
					ValidateOtp      bool       `inject:"config:core.auth.fake.validateOtp,optional"`
					UsernameFieldID  string     `inject:"config:core.auth.fake.usernameFieldId,optional"`
					PasswordFieldID  string     `inject:"config:core.auth.fake.passwordFieldId,optional"`
					OtpFieldID       string     `inject:"config:core.auth.fake.otpFieldId,optional"`
				})(testConfig),
			)

			ctx := web.ContextWithSession(context.Background(), web.EmptySession())
			got := c.Auth(ctx, tt.args.r)

			assert.Equal(t, tt.want, got)
		})
	}
}

func replaceStandardFormIDs(content string) string {
	result := strings.Replace(content, "{{.UsernameID}}", defaultUserNameFieldID, -1)
	result = strings.Replace(result, "{{.PasswordID}}", defaultPasswordFieldID, -1)
	result = strings.Replace(result, "{{.OtpID}}", defaultOtpFieldID, -1)

	return result
}
