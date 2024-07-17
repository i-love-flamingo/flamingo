//go:build tracelog

package flamingo

import (
	"context"
)

type (
	// Logger defines the Flamingo trace logger interface
	Logger interface {
		WithContext(ctx context.Context) Logger

		Trace(args ...interface{})
		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})

		Tracef(log string, args ...interface{})
		Debugf(log string, args ...interface{})

		WithField(key LogKey, value interface{}) Logger
		WithFields(fields map[LogKey]interface{}) Logger

		Flush()
	}
)
