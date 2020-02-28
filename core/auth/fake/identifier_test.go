package fake

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"text/template"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

type mockRouter struct {
	broker string
}

var _ web.ReverseRouter = (*mockRouter)(nil)

// Relative mock action
func (m *mockRouter) Relative(_ string, _ map[string]string) (*url.URL, error) {
	panic("not implemented")
}

// Absolute mock action
func (m *mockRouter) Absolute(_ *web.Request, _ string, _ map[string]string) (*url.URL, error) {
	return url.Parse(strings.ReplaceAll("/core/auth/login/:broker", ":broker", m.broker))
}

const testBroker = "testBroker"

func testIdentifier(config fakeConfig) *identifier {
	i := &identifier{
		config:        config,
		broker:        testBroker,
		responder:     new(web.Responder),
		eventRouter:   new(flamingo.DefaultEventRouter),
		reverseRouter: &mockRouter{broker: testBroker},
	}
	i.config.UsernameFieldID = "username"
	i.config.PasswordFieldID = "password"
	return i
}

func Test_Identifier_Authenticate(t *testing.T) {
	request := buildRequest(http.MethodGet, nil)
	testResponse := wantFormResponseWithMessage("")
	response := testIdentifier(fakeConfig{}).Authenticate(context.Background(), request)
	testResponse(t, response)
}

func Test_Identifier_Callback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		target       func(*identifier, context.Context, *web.Request) web.Result
		config       fakeConfig
		request      *web.Request
		testResponse func(*testing.T, web.Result)
		returnTo     func(*web.Request) *url.URL
	}{
		{
			name:         "simple identifier call with empty template to test fallback template",
			request:      buildRequest(http.MethodGet, nil),
			testResponse: wantFormResponseWithMessage(errIdentityNotSavedInSession.Error()),
		},
		{
			name:         "render auth form on empty form submit",
			request:      buildRequest(http.MethodPost, &url.Values{}),
			testResponse: wantFormResponseWithMessage(errMissingUsername.Error()),
		},
		{
			name: "renders auth form with message on invalid form data",
			request: buildRequest(http.MethodPost, &url.Values{
				"invalid-field": []string{"just to get into the form handling"},
			}),
			testResponse: wantFormResponseWithMessage(errMissingUsername.Error()),
		},
		{
			name: "show error message on invalid user",
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
			request: buildRequest(http.MethodPost, &url.Values{
				"username": []string{"nonexistent user"},
			}),
			testResponse: wantFormResponseWithMessage(errInvalidUser.Error()),
		},
		{
			name: "login for valid user without pwd",
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
			request: buildRequest(http.MethodPost, &url.Values{
				"username": []string{"user_b"},
			}),
			testResponse: func(t *testing.T, result web.Result) {
				response, ok := result.(*web.URLRedirectResponse)
				assert.True(t, ok)
				assert.NotNil(t, response.URL)
				assert.Equal(t, "/return/to/path", response.URL.Path)
			},
			returnTo: func(request *web.Request) *url.URL {
				return &url.URL{Path: "/return/to/path"}
			},
		},
		{
			name: "login for valid user with password mismatch",
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
			request: buildRequest(http.MethodPost, &url.Values{
				"username": []string{"user_b"},
				"password": []string{"invalid password"},
			}),
			testResponse: wantFormResponseWithMessage(errPasswordMismatch.Error()),
		},
		{
			name: "login for valid user / valid password",
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
			request: buildRequest(http.MethodPost, &url.Values{
				"username": []string{"user_b"},
				"password": []string{"test_b"},
			}),
			testResponse: func(t *testing.T, result web.Result) {
				response, ok := result.(*web.URLRedirectResponse)
				assert.True(t, ok)
				assert.NotNil(t, response.URL)
				assert.Equal(t, "/return/to/path", response.URL.Path)
			},
			returnTo: func(request *web.Request) *url.URL {
				return &url.URL{Path: "/return/to/path"}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := testIdentifier(tt.config)
			response := i.Callback(context.Background(), tt.request, tt.returnTo)
			tt.testResponse(t, response)
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

	result.Params = web.RequestParams{"broker": testBroker}

	return result
}

func wantFormResponseWithMessage(message string) func(*testing.T, web.Result) {
	return func(t *testing.T, result web.Result) {
		response, ok := result.(*web.Response)
		if !ok {
			t.Error("result is not of type *web.Response")
		}

		if response.Status != http.StatusOK {
			t.Errorf("status %d != 200", response.Status)
		}

		if contenttype := response.Header.Get("Content-Type"); contenttype != "text/html; charset=utf-8" {
			t.Errorf("Content-Type %q is not %q", contenttype, "text/html; charset=utf-8")
		}

		tpl := template.New("fake")
		tpl, err := tpl.Parse(defaultLoginTemplate)
		if err != nil {
			t.Error(err)
		}

		var body = new(bytes.Buffer)
		err = tpl.Execute(
			body,
			viewData{
				FormURL:    "/core/auth/login/testBroker",
				Message:    message,
				UsernameID: defaultUserNameFieldID,
				PasswordID: defaultPasswordFieldID,
			},
		)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, body, response.Body)
	}
}
