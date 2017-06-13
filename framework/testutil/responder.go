package testutil

import (
	"flamingo/framework/web"
)

type (
	// MockRedirectAware mock stub
	MockRedirectAware struct {
		CbRedirect             func(name string, args map[string]string) web.Response
		CbRedirectURL          func(url string) web.Response
		CbRedirectPermanent    func(name string, args map[string]string) web.Response
		CbRedirectPermanentURL func(url string) web.Response
	}

	// MockRenderAware mock stub
	MockRenderAware struct {
		CbRender func(context web.Context, tpl string, data interface{}) web.Response
	}

	// MockErrorAware mock stub
	MockErrorAware struct {
		CbError         func(context web.Context, err error) web.Response
		CbErrorNotFound func(context web.Context, error error) web.Response
	}
)

// Redirect mock
func (m *MockRedirectAware) Redirect(name string, args map[string]string) web.Response {
	return m.CbRedirect(name, args)
}

// RedirectURL mock
func (m *MockRedirectAware) RedirectURL(url string) web.Response {
	return m.CbRedirectURL(url)
}

// RedirectPermanent mock
func (m *MockRedirectAware) RedirectPermanent(name string, args map[string]string) web.Response {
	return m.CbRedirectPermanent(name, args)
}

// RedirectPermanentURL mock
func (m *MockRedirectAware) RedirectPermanentURL(url string) web.Response {
	return m.CbRedirectPermanentURL(url)
}

// Render mock
func (m *MockRenderAware) Render(context web.Context, tpl string, data interface{}) web.Response {
	return m.CbRender(context, tpl, data)
}

// Error mock
func (m *MockErrorAware) Error(context web.Context, err error) web.Response {
	return m.CbError(context, err)
}

// ErrorNotFound mock
func (m *MockErrorAware) ErrorNotFound(context web.Context, err error) web.Response {
	return m.CbErrorNotFound(context, err)
}
