//+build !go1.13

package flamingo

import (
	"fmt"
	"strings"
)

func init() {
	fmtErrorf = func(format string, a ...interface{}) error {
		return fmt.Errorf(strings.Replace(format, "%w", "%v", 1), a...)
	}
}
