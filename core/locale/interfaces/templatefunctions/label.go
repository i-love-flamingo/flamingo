package templatefunctions

import (
	"strings"
	//"github.com/nicksnyder/go-i18n"
)

type (
	// GetProduct is exported as a template function
	Label struct {
	}
)

// Name alias for use in template
func (tf Label) Name() string {
	return "__"
}

func (tf Label) Func() interface{} {
	return func(s ...string) string { return strings.Join(s, "::") }
}
