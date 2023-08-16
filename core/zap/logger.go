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

type Option func(*Logger)

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
		fields     []zapcore.Field
		logSession bool
	}
)

func NewLogger(logger *zap.Logger, options ...Option) *Logger {
	l := &Logger{
		Logger: logger,
	}

	for _, option := range options {
		option(l)
	}

	return l
}

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
	l.writeLog((*zap.Logger).Debug, fmt.Sprint(args...))
}

// Debugf logs a message at debug level with format string
func (l *Logger) Debugf(log string, args ...interface{}) {
	l.record("Debug")
	l.writeLog((*zap.Logger).Debug, fmt.Sprintf(log, args...))
}

// Info logs a message at info level
func (l *Logger) Info(args ...interface{}) {
	l.record("Info")
	l.writeLog((*zap.Logger).Info, fmt.Sprint(args...))
}

// Warn logs a message at warn level
func (l *Logger) Warn(args ...interface{}) {
	l.record("Warn")
	l.writeLog((*zap.Logger).Warn, fmt.Sprint(args...))
}

// Error logs a message at error level
func (l *Logger) Error(args ...interface{}) {
	l.record("Error")
	l.writeLog((*zap.Logger).Error, fmt.Sprint(args...))
}

// Fatal logs a message at fatal level
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func (l *Logger) Fatal(args ...interface{}) {
	l.record("Fatal")
	l.writeLog((*zap.Logger).Fatal, fmt.Sprint(args...))
}

// Panic logs a message at panic level
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panic(args ...interface{}) {
	l.record("Panic")
	l.writeLog((*zap.Logger).Panic, fmt.Sprint(args...))
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

	fields := make([]zapcore.Field, len(l.fields)+1)
	copy(fields, l.fields)
	fields[len(fields)-1] = zap.Any(string(key), value)

	return &Logger{
		Logger:     l.Logger.With(),
		configArea: area,
		fieldMap:   l.fieldMap,
		fields:     fields,
		logSession: l.logSession,
	}
}

// WithFields creates a child logger and adds structured context to it.
func (l *Logger) WithFields(fields map[flamingo.LogKey]interface{}) flamingo.Logger {
	newFields := make([]zap.Field, len(l.fields), len(l.fields)+len(fields))
	copy(newFields, l.fields)

	area := l.configArea

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

		newFields = append(newFields, zap.Any(string(key), value))
	}

	return &Logger{
		Logger:     l.Logger,
		configArea: area,
		fieldMap:   l.fieldMap,
		fields:     newFields,
		logSession: l.logSession,
	}
}

func (l *Logger) writeLog(logFunc func(zl *zap.Logger, msg string, fields ...zapcore.Field), msg string) {
	logFunc(l.Logger, msg, l.fields...)
}

// Flush is used by buffered loggers and triggers the actual writing. It is a good habit to call Flush before
// letting the process exit. For the top level flamingo.Logger, this is called by the app itself.
func (l *Logger) Flush() {
	_ = l.Logger.Sync()
}
