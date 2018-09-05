package flamingo

import (
	"context"
	"log"
)

//go:generate mockery -name "Logger"

// Common logger field keys
const (
	LogKeyAccesslog         LogKey = "accesslog" // LogKeyAccesslog marks a logmessage belonging to an (incoming) call (value should be 1)
	LogKeyApicall                  = "apicall"   // LogKeyApicall marks a logmessage belonging to an (outgoing) api call (value should be 1)
	LogKeyArea                     = "area"
	LogKeyBusinessID               = "businessId"
	LogKeyCategory                 = "category"
	LogKeySubCategory              = "sub_category"
	LogKeyClientIP                 = "client_ip"
	LogKeyCode                     = "code"
	LogKeyConnectionStatus         = "connection_status"
	LogKeyCorrelationID            = "correlationId"
	LogKeyLevel                    = "level"
	LogKeyMessage                  = "message"
	LogKeyMethod                   = "method"
	LogKeyPath                     = "path"
	LogKeyReferer                  = "referer"
	LogKeyRequest                  = "request"
	LogKeyRequestTime              = "request_time"
	LogKeyRequestedEndpoint        = "requested_endpoint"
	LogKeyRequestedURL             = "requested_url"
	LogKeyResponse                 = "response"
	LogKeyResponseCode             = "response_code"
	LogKeyResponseTime             = "response_time"
	LogKeySource                   = "source"
	LogKeyTimestamp                = "@timestamp"
	LogKeyTrace                    = "trace"
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

type StdLogger struct {
	log.Logger
}

func (l *StdLogger) Debug(args ...interface{}) {
	l.Print(args...)
}

func (l *StdLogger) Debugf(f string, args ...interface{}) {
	l.Printf(f, args...)
}

func (l *StdLogger) Info(args ...interface{}) {
	l.Print(args...)
}

func (l *StdLogger) Warn(args ...interface{}) {
	l.Print(args...)
}

func (l *StdLogger) WithContext(ctx context.Context) Logger {
	return l
}

func (l *StdLogger) WithField(key LogKey, value interface{}) Logger {
	log.Println("WithField", key, value)
	return l
}

func (l *StdLogger) WithFields(fields map[LogKey]interface{}) Logger {
	log.Println("WithFields", fields)
	return l
}

func (l *StdLogger) Error(args ...interface{}) {
	l.Print(args...)
}

func (l *StdLogger) Flush() {}

// NullLogger does not log
type NullLogger struct{}

// WithContext null-implementation
func (n NullLogger) WithContext(ctx context.Context) Logger { return n }

// WithField null-implementation
func (n NullLogger) WithField(key LogKey, value interface{}) Logger { return n }

// WithFields null-implementation
func (n NullLogger) WithFields(fields map[LogKey]interface{}) Logger { return n }

// Debug null-implementation
func (NullLogger) Debug(args ...interface{}) {}

// Debug null-implementation
func (NullLogger) Debugf(log string, args ...interface{}) {}

// Info null-implementation
func (NullLogger) Info(args ...interface{}) {}

// Warn null-implementation
func (NullLogger) Warn(args ...interface{}) {}

// Error null-implementation
func (NullLogger) Error(args ...interface{}) {}

// Fatal null-implementation
func (NullLogger) Fatal(args ...interface{}) {}

// Panic null-implementation
func (NullLogger) Panic(args ...interface{}) {}

// Flush null-implementation
func (n NullLogger) Flush() {}
