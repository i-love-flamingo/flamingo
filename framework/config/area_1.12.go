//+build go1.11 go1.12

package config

import (
	"fmt"
	"strings"
)

func init() {
	fmtErrorf = func(format string, a ...interface{}) error {
		return fmt.Errorf(strings.Replace(format, "%w", "%v", 1), a...)
	}
}
