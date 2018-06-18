package domain

import (
	"fmt"

	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger is a Wrapper for the zap logger fulfilling the flamingo.Logger interface
	Logger struct {
		*zap.Logger
	}
)

// WithCorrelationID returns a logger with a correlation ID field
func (l *Logger) WithCorrelationID(cid string) flamingo.Logger {
	return l.WithField("correlationId", cid)
}

// WithContext returns a logger with fields filled from the context
// businessId:    From Header X-Business-ID
// client_ip:     From Header X-Forwarded-For or request if header is empty
// correlationId: The ID of the context
// method:        HTTP verb from request
// path:          URL path from request
// referer:       referer from request
// request:       received payload from request
func (l *Logger) WithContext(ctx web.Context) flamingo.Logger {
	request := ctx.Request()
	clientIP := request.RemoteAddr
	if request.Header.Get("X-Forwarded-For") != "" {
		clientIP += ", " + request.Header.Get("X-Forwarded-For")
	}
	//body, _ := ioutil.ReadAll(request.Body)

	return l.WithFields(
		map[flamingo.LogKey]interface{}{
			flamingo.LogKeyBusinessID:    request.Header.Get("X-Business-ID"),
			flamingo.LogKeyClientIP:      clientIP,
			flamingo.LogKeyCorrelationID: ctx.ID(),
			flamingo.LogKeyMethod:        request.Method,
			flamingo.LogKeyPath:          request.URL.Path,
			flamingo.LogKeyReferer:       request.Referer(),
			//flamingo.LogKeyRequest:       string(body),
		},
	)
}

// Debug logs a message at debug level
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(fmt.Sprint(args...))
}

// Info logs a message at info level
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(fmt.Sprint(args...))
}

// Warn logs a message at warn level
func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(fmt.Sprint(args...))
}

// Error logs a message at error level
func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(fmt.Sprint(args...))
}

// Fatal logs a message at fatal level
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(fmt.Sprint(args...))
}

// Panic logs a message at panic level
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panic(args ...interface{}) {
	l.Logger.Panic(fmt.Sprint(args...))
}

// WithField creates a child logger and adds structured context to it.
func (l *Logger) WithField(key flamingo.LogKey, value interface{}) flamingo.Logger {
	return &Logger{
		Logger: l.Logger.With(
			zap.Any(string(key), value),
		),
	}
}

// WithFields creates a child logger and adds structured context to it.
func (l *Logger) WithFields(fields map[flamingo.LogKey]interface{}) flamingo.Logger {
	zapFields := make([]zapcore.Field, len(fields))
	i := 0
	for key, value := range fields {
		zapFields[i] = zap.Any(string(key), value)
		i++
	}

	return &Logger{
		Logger: l.Logger.With(zapFields...),
	}
}

// Flush is used by buffered loggers and triggers the actual writing. It is a good habit to call Flush before
// letting the process exit. For the top level flamingo.Logger, this is called by the app itself.
func (l *Logger) Flush() {
	l.Logger.Sync()
}
