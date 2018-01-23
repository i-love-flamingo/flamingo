package templatefunctions

import (
	"strings"
)

type (
	// startsWithFunc returns boolean
	StartsWithFunc struct{}
)

// Name alias for use in template
func (s StartsWithFunc) Name() string {
	return "startsWith"
}

// Func as implementation of url method
func (s *StartsWithFunc) Func() interface{} {
	return func(haystack string, needle string) bool {
		haystack = strings.ToLower(haystack)
		needle = strings.ToLower(needle)
		return strings.HasPrefix(haystack, needle)
	}
}
