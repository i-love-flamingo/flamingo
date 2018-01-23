package templatefunctions

import (
	"strings"
)

type (
	StartsWithFunc struct{}
)

func (s StartsWithFunc) Name() string {
	return "startsWith"
}

func (s *StartsWithFunc) Func() interface{} {
	return func(haystack string, needle string) bool {
		haystack = strings.ToLower(haystack)
		needle = strings.ToLower(needle)
		return strings.HasPrefix(haystack, needle)
	}
}
