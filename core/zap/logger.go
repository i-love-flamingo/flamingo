package zap

import (
	"context"
	"fmt"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
)

var (
	logCount    = stats.Int64("flamingo/zap/logs", "Count of logs", stats.UnitDimensionless)
	keyLevel, _ = tag.NewKey("level")
)

func init() {
	if err := opencensus.View("flamingo/zap/logs", logCount, view.Count(), keyLevel); err != nil {
		panic(err)
	}
}

type (
	// Logger is a Wrapper for the zap logger fulfilling the flamingo.Logger interface
	Logger struct {
		*zap.Logger
		configArea string
		fieldMap   map[string]string
		logSession bool
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

func (l *Logger) record(level string) {
	if !l.Core().Enabled(logLevels[level]) {
		return
	}

	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, l.configArea), tag.Upsert(keyLevel, level))
	stats.Record(ctx, logCount.M(1))
}

// Debug logs a message at debug level
func (l *Logger) Debug(args ...interface{}) {
	l.record("Debug")
	l.Logger.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at debug level with format string
func (l *Logger) Debugf(log string, args ...interface{}) {
	l.record("Debug")
	l.Logger.Debug(fmt.Sprintf(log, args...))
}

// Info logs a message at info level
func (l *Logger) Info(args ...interface{}) {
	l.record("Info")
	l.Logger.Info(fmt.Sprint(args...))
}

// Warn logs a message at warn level
func (l *Logger) Warn(args ...interface{}) {
	l.record("Warn")
	l.Logger.Warn(fmt.Sprint(args...))
}

// Error logs a message at error level
func (l *Logger) Error(args ...interface{}) {
	l.record("Error")
	l.Logger.Error(fmt.Sprint(args...))
}

// Fatal logs a message at fatal level
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func (l *Logger) Fatal(args ...interface{}) {
	l.record("Fatal")
	l.Logger.Fatal(fmt.Sprint(args...))
}

// Panic logs a message at panic level
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panic(args ...interface{}) {
	l.record("Panic")
	l.Logger.Panic(fmt.Sprint(args...))
}

// WithField creates a child logger and adds structured context to it.
func (l *Logger) WithField(key flamingo.LogKey, value interface{}) flamingo.Logger {
	area := l.configArea
	if key == flamingo.LogKeyArea {
		area = value.(string)
	}

	if alias, ok := l.fieldMap[string(key)]; ok {
		// skip field
		if alias == "-" {
			return l
		}

		key = flamingo.LogKey(alias)
	}

	return &Logger{
		Logger:     l.Logger.With(zap.Any(string(key), value)),
		configArea: area,
		fieldMap:   l.fieldMap,
		logSession: l.logSession,
	}
}

// WithFields creates a child logger and adds structured context to it.
func (l *Logger) WithFields(fields map[flamingo.LogKey]interface{}) flamingo.Logger {
	zapFields := make([]zapcore.Field, len(fields))
	area := l.configArea
	i := 0
	for key, value := range fields {
		if key == flamingo.LogKeyArea {
			area = value.(string)
		}

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
		configArea: area,
		fieldMap:   l.fieldMap,
		logSession: l.logSession,
	}
}

// Flush is used by buffered loggers and triggers the actual writing. It is a good habit to call Flush before
// letting the process exit. For the top level flamingo.Logger, this is called by the app itself.
func (l *Logger) Flush() {
	l.Logger.Sync()
}
