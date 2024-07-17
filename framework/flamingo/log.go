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
)

func LogFunc(l Logger, fields map[LogKey]any) func(f func(l Logger, args ...any), args ...any) {
	if l == nil {
		l = new(NullLogger)
	}

	return func(f func(n Logger, args ...any), args ...any) {
		if f == nil {
			f = Logger.Info
		}

		f(l.WithFields(fields), args...)
	}
}

func LogFuncWithContext(l Logger, fields map[LogKey]any) func(ctx context.Context, f func(l Logger, args ...any), args ...any) {
	if l == nil {
		l = new(NullLogger)
	}

	return func(ctx context.Context, f func(n Logger, args ...any), args ...any) {
		if f == nil {
			f = Logger.Info
		}

		f(l.WithContext(ctx).WithFields(fields), args...)
	}
}

var _ Logger = new(NullLogger)
var _ Logger = new(StdLogger)

// StdLogger uses the go stdlib logger for logging
type StdLogger struct {
	Logger log.Logger
}

func isLoggerNil(l *StdLogger) bool {
	return l == nil || l.Logger.Writer() == nil
}

func (l *StdLogger) Fatal(args ...interface{}) {
	if isLoggerNil(l) {
		log.Fatalln(args...)
		return
	}

	l.Logger.Fatal(args...)
}

func (l *StdLogger) Panic(args ...interface{}) {
	if isLoggerNil(l) {
		log.Panic(args...)
		return
	}

	l.Logger.Panic(args...)
}

// Trace logs output
func (l *StdLogger) Trace(args ...interface{}) {
	args = append([]any{"trace: "}, args...)
	if isLoggerNil(l) {
		log.Print(args...)
		return
	}

	l.Logger.Print(args...)
}

// Tracef outputs the formatted trace string
func (l *StdLogger) Tracef(f string, args ...interface{}) {
	if isLoggerNil(l) {
		log.Print(args...)
		return
	}

	l.Logger.Printf(f, args...)
}

// Debug logs output
func (l *StdLogger) Debug(args ...interface{}) {
	args = append([]any{"debug: "}, args...)
	if isLoggerNil(l) {
		log.Print(args...)
		return
	}

	l.Logger.Print(args...)
}

// Debugf outputs the formatted debug string
func (l *StdLogger) Debugf(f string, args ...interface{}) {
	if isLoggerNil(l) {
		log.Print(args...)
		return
	}

	l.Logger.Printf(f, args...)
}

// Info log output
func (l *StdLogger) Info(args ...interface{}) {
	args = append([]any{"info: "}, args...)
	if isLoggerNil(l) {
		log.Print(args...)
		return
	}

	l.Logger.Print(args...)
}

// Warn log output
func (l *StdLogger) Warn(args ...interface{}) {
	args = append([]any{"warn: "}, args...)
	if isLoggerNil(l) {
		log.Print(args...)
		return
	}

	l.Logger.Print(args...)
}

// WithContext currently does nothing
func (l *StdLogger) WithContext(_ context.Context) Logger {
	return l
}

// WithField currently logs the field
func (l *StdLogger) WithField(key LogKey, value interface{}) Logger {
	if isLoggerNil(l) {
		log.Println("WithField", key, value)
		return l
	}

	l.Logger.Println("WithField", key, value)

	return l
}

// WithFields currently logs the fields
func (l *StdLogger) WithFields(fields map[LogKey]interface{}) Logger {
	if isLoggerNil(l) {
		log.Println("WithFields", fields)
		return l
	}

	l.Logger.Println("WithFields", fields)

	return l
}

// Error log
func (l *StdLogger) Error(args ...interface{}) {
	args = append([]any{"error: "}, args...)
	if isLoggerNil(l) {
		log.Print(args...)
		return
	}

	l.Logger.Print(args...)
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

// Trace null-implementation
func (NullLogger) Trace(_ ...interface{}) {}

// Tracef null-implementation
func (NullLogger) Tracef(_ string, _ ...interface{}) {}

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
