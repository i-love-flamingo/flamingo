package web

import (
	"context"
	"strings"
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
func ToContext(ctx context.Context) Context {
	return ctx.Value(CONTEXT).(Context)
}

// ToRequest upgrades a web.Context to the new context+request form
func ToRequest(ctx Context) (context.Context, *Request) {
	request := ctx.Value("__req").(*Request)

	return request.request.Context(), request
}
