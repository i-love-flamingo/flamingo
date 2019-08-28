package zap

import (
	"context"
	"flamingo.me/flamingo/v3/core/zap/application"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"fmt"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger is a Wrapper for the zap logger fulfilling the flamingo.Logger interface
	Logger struct {
		*zap.Logger
		fieldMap        map[string]string
		logSession      bool
	}
)

// WithContext returns a logger with fields filled from the context
// businessId:    From Header X-Business-ID
// client_ip:     From Header X-Forwarded-For or request if header is empty
// correlationId: The ID of the context
// method:        HTTP verb from request
// path:          URL path from request
// referer:       referer from request
// request:       received payload from request
func (l *Logger) WithContext(ctx context.Context) flamingo.Logger {
	span := trace.FromContext(ctx)
	fields := map[flamingo.LogKey]interface{}{
		flamingo.LogKeyTraceID: span.SpanContext().TraceID.String(),
		flamingo.LogKeySpanID:  span.SpanContext().SpanID.String(),
	}

	req := web.RequestFromContext(ctx)

	if req != nil {
		request := req.Request()
		fields[flamingo.LogKeyMethod] = request.Method
		fields[flamingo.LogKeyPath] = request.URL.Path
	}

	if l.logSession {
		session := web.SessionFromContext(ctx)
		if session != nil {
			fields[flamingo.LogKeySession] = session.IDHash()
		}
	}
	return l.WithFields(fields)
}

// Debug logs a message at debug level
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at debug level with format string
func (l *Logger) Debugf(log string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(log, args...))
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

	go func() {
		ctx, _ := tag.New(
			context.Background(),
			tag.Update(opencensus.KeyArea, "root"),
		)
		stats.Record(ctx, application.ErrorCount.M(1))
	}()
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
	if alias, ok := l.fieldMap[string(key)]; ok {
		// skip field
		if alias == "-" {
			return l
		}

		key = flamingo.LogKey(alias)
	}

	return &Logger{
		Logger:     l.Logger.With(zap.Any(string(key), value)),
		fieldMap:   l.fieldMap,
		logSession: l.logSession,
	}
}

// WithFields creates a child logger and adds structured context to it.
func (l *Logger) WithFields(fields map[flamingo.LogKey]interface{}) flamingo.Logger {
	zapFields := make([]zapcore.Field, len(fields))
	i := 0
	for key, value := range fields {
		if alias, ok := l.fieldMap[string(key)]; ok {
			// skip field
			if alias == "-" {
				continue
			}

			key = flamingo.LogKey(alias)
		}

		zapFields[i] = zap.Any(string(key), value)
		i++
	}

	return &Logger{
		Logger:     l.Logger.With(zapFields[:i]...),
		fieldMap:   l.fieldMap,
		logSession: l.logSession,
	}
}

// Flush is used by buffered loggers and triggers the actual writing. It is a good habit to call Flush before
// letting the process exit. For the top level flamingo.Logger, this is called by the app itself.
func (l *Logger) Flush() {
	l.Logger.Sync()
}
