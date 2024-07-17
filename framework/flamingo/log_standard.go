//go:build !tracelog

package flamingo

import (
	"context"
)

type (
	// Logger defines the standard Flamingo logger interface
	Logger interface {
		WithContext(ctx context.Context) Logger

		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})

		Debugf(log string, args ...interface{})

		WithField(key LogKey, value interface{}) Logger
		WithFields(fields map[LogKey]interface{}) Logger

		Flush()
	}
)
