package templatefunctions

import "strings"

type (
	CapitalizeFunc struct{}
)

func (s *CapitalizeFunc) Func() interface{} {
	return func(str string) string {
		return strings.Title(str)
	}
}
