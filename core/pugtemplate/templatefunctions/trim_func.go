package templatefunctions

import "strings"

type (
	TrimFunc struct{}
)

func (s *TrimFunc) Func() interface{} {
	return func(str string) string {
		return strings.TrimSpace(str)
	}
}
