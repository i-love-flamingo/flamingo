package flamingo

import (
	"go.aoe.com/flamingo/framework/web"
)

//go:generate mockery -name "Logger"

// Common logger field keys
const (
	// LogKeyAccesslog marks a logmessage belonging to an (incoming) call (value should be 1)
	LogKeyAccesslog = "accesslog"
	// LogKeyApicall marks a logmessage belonging to an (outgoing) api call (value should be 1)
	LogKeyApicall           = "apicall"
	LogKeyArea              = "area"
	LogKeyBusinessID        = "businessId"
	LogKeyCategory          = "category"
	LogKeyClientIP          = "client_ip"
	LogKeyCode              = "code"
	LogKeyConnectionStatus  = "connection_status"
	LogKeyCorrelationID     = "correlationId"
	LogKeyLevel             = "level"
	LogKeyMessage           = "message"
	LogKeyMethod            = "method"
	LogKeyPath              = "path"
	LogKeyReferer           = "referer"
	LogKeyRequest           = "request"
	LogKeyRequestTime       = "request_time"
	LogKeyRequestedEndpoint = "requested_endpoint"
	LogKeyRequestedURL      = "requested_url"
	LogKeyResponse          = "response"
	LogKeyResponseCode      = "response_code"
	LogKeyResponseTime      = "response_time"
	LogKeySource            = "source"
	LogKeyTimestamp         = "@timestamp"
	LogKeyTrace             = "trace"
)

type (
	// Logger defines a standard Flamingo logger interfaces
	Logger interface {
		WithContext(ctx web.Context) Logger

		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})

		WithField(key string, value interface{}) Logger
		WithFields(fields map[string]interface{}) Logger

		Flush()
	}
)

// NullLogger does not log
type NullLogger struct{}

func (n NullLogger) WithContext(ctx web.Context) Logger              { return n }
func (n NullLogger) WithField(key string, value interface{}) Logger  { return n }
func (n NullLogger) WithFields(fields map[string]interface{}) Logger { return n }
func (NullLogger) Debug(args ...interface{})                         {}
func (NullLogger) Info(args ...interface{})                          {}
func (NullLogger) Print(args ...interface{})                         {}
func (NullLogger) Warn(args ...interface{})                          {}
func (NullLogger) Error(args ...interface{})                         {}
func (NullLogger) Fatal(args ...interface{})                         {}
func (NullLogger) Panic(args ...interface{})                         {}
func (n NullLogger) Flush()                                          {}
