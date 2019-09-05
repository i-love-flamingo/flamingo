package web

import (
	"context"
	"net/url"
)

type (
	partialDataContextKey string

	// SetPartialDataFunc allows to set partial data
	SetPartialDataFunc struct{}

	// GetPartialDataFunc allows to get partial data
	GetPartialDataFunc struct{}
)

const ctxKey partialDataContextKey = "partialData"

// Func getter to bind the context
func (*SetPartialDataFunc) Func(c context.Context) interface{} {
	return func(key string, val interface{}) interface{} {
		r := RequestFromContext(c)
		if r == nil {
			return nil
		}

		data, ok := r.Values.Load(ctxKey)
		if !ok || data == nil {
			data = make(map[string]interface{})
		}

		data.(map[string]interface{})[key] = val

		r.Values.Store(ctxKey, data)

		return nil
	}
}

// Func getter to bind the context
func (*GetPartialDataFunc) Func(c context.Context) interface{} {
	return func() map[string]interface{} {
		r := RequestFromContext(c)
		if r == nil {
			return nil
		}

		data, ok := r.Values.Load(ctxKey)
		if !ok || data == nil {
			return nil
		}

		return data.(map[string]interface{})
	}
}

// CanonicalDomainFunc is exported as a template function
type CanonicalDomainFunc struct {
	router ReverseRouter
}

// Inject dependencies
func (c *CanonicalDomainFunc) Inject(router ReverseRouter) *CanonicalDomainFunc {
	c.router = router
	return c
}

// Func returns the canonicalDomain func
func (c *CanonicalDomainFunc) Func(ctx context.Context) interface{} {
	return func() string {
		u, _ := c.router.Absolute(RequestFromContext(ctx), "", nil)
		return u.String()
	}
}

// IsExternalURL is exported as a template function
type IsExternalURL struct {
	router ReverseRouter
}

// Inject dependencies
func (c *IsExternalURL) Inject(router ReverseRouter) *IsExternalURL {
	c.router = router
	return c
}

// Func returns a boolean if a given URL is external
func (c *IsExternalURL) Func(ctx context.Context) interface{} {
	return func(urlStr string) bool {
		if u, err := url.Parse(urlStr); err == nil {
			au, _ := c.router.Absolute(RequestFromContext(ctx), "", nil)
			return u.Host != "" && au.Host != u.Host
		}

		return false
	}
}
