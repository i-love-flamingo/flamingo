package web

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

// URLTitle normalizes a title for nice usage in URLs
func URLTitle(title string) string {
	url := strings.ToLower(strings.Replace(strings.Replace(title, "/", "_", -1), " ", "-", -1))
	url = strings.Replace(url, "-_", "-", -1)
	url = strings.Replace(url, "%", "-", -1)
	for strings.Contains(url, "--") {
		url = strings.Replace(url, "--", "-", -1)
	}

	return url
}

// ToContext supports legacy usage of the web.Context type
// deprecated: always use context+request directly
func ToContext(ctx context.Context) Context {
	if c, ok := ctx.(Context); ok {
		return c
	}
	if c, ok := ctx.Value(CONTEXT).(Context); ok {
		return c
	}

	panic(errors.New("can not convert a context.Context to a web.Context"))
}

// ToRequest upgrades a web.Context to the new context+request form
// deprecated: always use context+request directly
func ToRequest(ctx Context) (context.Context, *Request) {
	request, ok := ctx.Value("__req").(*Request)

	if !ok {
		return ctx, RequestFromRequest(ctx.Request(), ctx.Session())
	}

	return request.request.Context(), request
}
