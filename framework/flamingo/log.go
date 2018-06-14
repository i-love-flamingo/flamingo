package flamingo

import (
	"go.aoe.com/flamingo/framework/web"
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
	LogKey string

	// Logger defines a standard Flamingo logger interfaces
	Logger interface {
		WithContext(ctx web.Context) Logger

		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})

		WithField(key LogKey, value interface{}) Logger
		WithFields(fields map[LogKey]interface{}) Logger

		Flush()
	}
)

// NullLogger does not log
type NullLogger struct{}

func (n NullLogger) WithContext(ctx web.Context) Logger              { return n }
func (n NullLogger) WithField(key LogKey, value interface{}) Logger  { return n }
func (n NullLogger) WithFields(fields map[LogKey]interface{}) Logger { return n }
func (NullLogger) Debug(args ...interface{})                         {}
func (NullLogger) Info(args ...interface{})                          {}
func (NullLogger) Print(args ...interface{})                         {}
func (NullLogger) Warn(args ...interface{})                          {}
func (NullLogger) Error(args ...interface{})                         {}
func (NullLogger) Fatal(args ...interface{})                         {}
func (NullLogger) Panic(args ...interface{})                         {}
func (n NullLogger) Flush()                                          {}
