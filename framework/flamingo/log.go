package flamingo

import (
	"context"
	"log"
)

// Common logger field keys
const (
	LogKeyAccesslog         LogKey = "accesslog" // LogKeyAccesslog marks a logmessage belonging to an (incoming) call (value should be 1)
	LogKeyApicall           LogKey = "apicall"   // LogKeyApicall marks a logmessage belonging to an (outgoing) api call (value should be 1)
	LogKeyArea              LogKey = "area"
	LogKeyBusinessID        LogKey = "businessId"
	LogKeyCategory          LogKey = "category"
	LogKeyModule            LogKey = "module"
	LogKeySubCategory       LogKey = "sub_category"
	LogKeyClientIP          LogKey = "client_ip"
	LogKeyCode              LogKey = "code"
	LogKeyConnectionStatus  LogKey = "connection_status"
	LogKeyCorrelationID     LogKey = "correlationId"
	LogKeyTraceID           LogKey = "traceID"
	LogKeySpanID            LogKey = "spanID"
	LogKeyLevel             LogKey = "level"
	LogKeyMessage           LogKey = "message"
	LogKeyMethod            LogKey = "method"
	LogKeySession           LogKey = "session"
	LogKeyPath              LogKey = "path"
	LogKeyReferer           LogKey = "referer"
	LogKeyRequest           LogKey = "request"
	LogKeyRequestTime       LogKey = "request_time"
	LogKeyRequestedEndpoint LogKey = "requested_endpoint"
	LogKeyRequestedURL      LogKey = "requested_url"
	LogKeyResponse          LogKey = "response"
	LogKeyResponseCode      LogKey = "response_code"
	LogKeyResponseTime      LogKey = "response_time"
	LogKeySource            LogKey = "source"
	LogKeyTimestamp         LogKey = "@timestamp"
	LogKeyTrace             LogKey = "trace"
)

type (
	// LogKey is a logging key constant
	LogKey string

	// Logger defines a standard Flamingo logger interfaces
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

var _ Logger = new(NullLogger)
var _ Logger = new(StdLogger)

// StdLogger uses the go stdlib logger for logging
type StdLogger struct {
	log.Logger
}

// Debug logs output
func (l *StdLogger) Debug(args ...interface{}) {
	l.Print(args...)
}

// Debugf outputs the formatted debug string
func (l *StdLogger) Debugf(f string, args ...interface{}) {
	l.Printf(f, args...)
}

// Info log output
func (l *StdLogger) Info(args ...interface{}) {
	l.Print(args...)
}

// Warn log output
func (l *StdLogger) Warn(args ...interface{}) {
	l.Print(args...)
}

// WithContext currently does nothing
func (l *StdLogger) WithContext(_ context.Context) Logger {
	return l
}

// WithField currently logs the field
func (l *StdLogger) WithField(key LogKey, value interface{}) Logger {
	log.Println("WithField", key, value)
	return l
}

// WithFields currently logs the fields
func (l *StdLogger) WithFields(fields map[LogKey]interface{}) Logger {
	log.Println("WithFields", fields)
	return l
}

// Error log
func (l *StdLogger) Error(args ...interface{}) {
	l.Print(args...)
}

// Flush does nothing
func (l *StdLogger) Flush() {}

// NullLogger does not log
type NullLogger struct{}

// WithContext null-implementation
func (n NullLogger) WithContext(_ context.Context) Logger { return n }

// WithField null-implementation
func (n NullLogger) WithField(_ LogKey, _ interface{}) Logger { return n }

// WithFields null-implementation
func (n NullLogger) WithFields(_ map[LogKey]interface{}) Logger { return n }

// Debug null-implementation
func (NullLogger) Debug(_ ...interface{}) {}

// Debugf null-implementation
func (NullLogger) Debugf(_ string, _ ...interface{}) {}

// Info null-implementation
func (NullLogger) Info(_ ...interface{}) {}

// Warn null-implementation
func (NullLogger) Warn(_ ...interface{}) {}

// Error null-implementation
func (NullLogger) Error(_ ...interface{}) {}

// Fatal null-implementation
func (NullLogger) Fatal(_ ...interface{}) {}

// Panic null-implementation
func (NullLogger) Panic(_ ...interface{}) {}

// Flush null-implementation
func (n NullLogger) Flush() {}
