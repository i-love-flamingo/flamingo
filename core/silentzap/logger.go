package silentzap

import (
	"context"
	"fmt"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// SilentLogger is a Wrapper for the zap logger fulfilling the flamingo.Logger interface
	SilentLogger struct {
		*zap.Logger
		configArea      string
		fieldMap        map[string]string
		logSession      bool
		loggingRegistry *LoggingContextRegistry
		traceID         string
	}
)

var (
	logCount    = stats.Int64("flamingo/silentlogger/logs", "Count of logs", stats.UnitDimensionless)
	keyLevel, _ = tag.NewKey("level")

	logLevels = map[string]zapcore.Level{
		"Debug":  zap.DebugLevel,
		"Info":   zap.InfoLevel,
		"Warn":   zap.WarnLevel,
		"Error":  zap.ErrorLevel,
		"DPanic": zap.DPanicLevel,
		"Panic":  zap.PanicLevel,
		"Fatal":  zap.FatalLevel,
	}
)

func init() {
	if err := opencensus.View("flamingo/silentlogger/logs", logCount, view.Count(), keyLevel); err != nil {
		panic(err)
	}
}

func getSilentLogger(
	loggingRegistry *LoggingContextRegistry,
	config *struct {
		Area               string     `inject:"config:area"`
		JSON               bool       `inject:"config:core.zap.json,optional"`
		LogLevel           string     `inject:"config:core.zap.loglevel,optional"`
		ColoredOutput      bool       `inject:"config:core.zap.colored,optional"`
		DevelopmentMode    bool       `inject:"config:core.zap.devmode,optional"`
		SamplingEnabled    bool       `inject:"config:core.zap.sampling.enabled,optional"`
		SamplingInitial    float64    `inject:"config:core.zap.sampling.initial,optional"`
		SamplingThereafter float64    `inject:"config:core.zap.sampling.thereafter,optional"`
		FieldMap           config.Map `inject:"config:core.zap.fieldmap,optional"`
		LogSession         bool       `inject:"config:core.zap.logsession,optional"`
	},
) *SilentLogger {
	level, ok := logLevels[config.LogLevel]
	if !ok {
		// if nothing is configured user ErrorLevel
		level = zap.ErrorLevel
	}

	var samplingConfig *zap.SamplingConfig
	if config.SamplingEnabled && config.SamplingThereafter > 0 && config.SamplingInitial > 0 {
		samplingConfig = &zap.SamplingConfig{
			Initial:    int(config.SamplingInitial),
			Thereafter: int(config.SamplingThereafter),
		}
	}

	output := "console"
	if config.JSON {
		output = "json"
	}

	encoder := zapcore.CapitalLevelEncoder
	if config.ColoredOutput {
		encoder = zapcore.CapitalColorLevelEncoder
	}

	cfg := makeZapConfig(level, config.DevelopmentMode, samplingConfig, output, encoder)

	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	fieldMap := makeFieldMap(config.FieldMap)

	silentLogger := &SilentLogger{
		Logger:          logger,
		fieldMap:        fieldMap,
		logSession:      config.LogSession,
		configArea:      config.Area,
		loggingRegistry: loggingRegistry,
	}

	silentLogger, ok = silentLogger.WithField(flamingo.LogKeyArea, config.Area).(*SilentLogger)
	if !ok {
		panic("getSilentLogger just tried to create logger of different type")
	}

	return silentLogger
}

func makeZapConfig(
	level zapcore.Level,
	developmentMode bool,
	samplingConfig *zap.SamplingConfig,
	encoding string,
	encoder zapcore.LevelEncoder,
) zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       developmentMode,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          samplingConfig,
		Encoding:          encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     string(flamingo.LogKeyMessage),
			LevelKey:       string(flamingo.LogKeyLevel),
			TimeKey:        string(flamingo.LogKeyTimestamp),
			NameKey:        "logger",
			CallerKey:      string(flamingo.LogKeySource),
			StacktraceKey:  string(flamingo.LogKeyTrace),
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    nil,
	}
}

func makeFieldMap(configFieldMap config.Map) map[string]string {
	fieldMap := make(map[string]string, len(configFieldMap))

	for k, v := range configFieldMap {
		if v, ok := v.(string); ok {
			fieldMap[k] = v
		}
	}

	return fieldMap
}

func (l *SilentLogger) record(level string) {
	if !l.Core().Enabled(logLevels[level]) {
		return
	}

	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, l.configArea), tag.Upsert(keyLevel, level))
	stats.Record(ctx, logCount.M(1))
}

// Debug logs a message at debug level
func (l *SilentLogger) Debug(args ...interface{}) {
	l.record("Debug")

	logContext := l.loggingRegistry.Get(l.traceID)
	if logContext.isWritingAllowed() {
		l.Logger.Debug(fmt.Sprint(args...))
		return
	}

	checkedEntry := l.Logger.Check(zapcore.DebugLevel, fmt.Sprint(args...))
	logContext.store(checkedEntry)
}

// Debugf logs a message at debug level with format string
func (l *SilentLogger) Debugf(log string, args ...interface{}) {
	l.record("Debug")

	logContext := l.loggingRegistry.Get(l.traceID)
	if logContext.isWritingAllowed() {
		l.Logger.Debug(fmt.Sprintf(log, args...))
		return
	}

	checkedEntry := l.Logger.Check(zapcore.DebugLevel, fmt.Sprintf(log, args...))
	logContext.store(checkedEntry)
}

// Info logs a message at info level
func (l *SilentLogger) Info(args ...interface{}) {
	l.record("Info")

	logContext := l.loggingRegistry.Get(l.traceID)
	if logContext.isWritingAllowed() {
		l.Logger.Info(fmt.Sprint(args...))
		return
	}

	checkedEntry := l.Logger.Check(zapcore.InfoLevel, fmt.Sprint(args...))
	logContext.store(checkedEntry)
}

// Warn logs a message at warn level
func (l *SilentLogger) Warn(args ...interface{}) {
	l.record("Warn")

	logContext := l.loggingRegistry.Get(l.traceID)
	if logContext.isWritingAllowed() {
		l.Logger.Warn(fmt.Sprint(args...))
		return
	}

	checkedEntry := l.Logger.Check(zapcore.WarnLevel, fmt.Sprint(args...))
	logContext.store(checkedEntry)
}

// Error logs a message at error level
func (l *SilentLogger) Error(args ...interface{}) {
	l.record("Error")
	logContext := l.loggingRegistry.Get(l.traceID)

	if !logContext.isWritingAllowed() {
		currentEntries := logContext.get()
		for _, entry := range currentEntries {
			entry.CheckedLogEntry.Write(entry.Fields...)
		}
	}

	l.Logger.Error(fmt.Sprint(args...))
}

// Fatal logs a message at fatal level
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func (l *SilentLogger) Fatal(args ...interface{}) {
	l.record("Fatal")

	logContext := l.loggingRegistry.Get(l.traceID)

	if !logContext.isWritingAllowed() {
		currentEntries := logContext.get()
		for _, entry := range currentEntries {
			entry.CheckedLogEntry.Write(entry.Fields...)
		}
	}

	l.Logger.Fatal(fmt.Sprint(args...))
}

// Panic logs a message at panic level
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *SilentLogger) Panic(args ...interface{}) {
	l.record("Panic")

	logContext := l.loggingRegistry.Get(l.traceID)

	if !logContext.isWritingAllowed() {
		currentEntries := logContext.get()
		for _, entry := range currentEntries {
			entry.CheckedLogEntry.Write(entry.Fields...)
		}
	}

	l.Logger.Panic(fmt.Sprint(args...))
}

func (l *SilentLogger) WithContext(ctx context.Context) flamingo.Logger {
	span := trace.FromContext(ctx)
	traceId := span.SpanContext().TraceID.String()
	fields := map[flamingo.LogKey]interface{}{
		flamingo.LogKeyTraceID: traceId,
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

// WithFields creates a child logger and adds structured context to it.
func (l *SilentLogger) WithFields(fields map[flamingo.LogKey]interface{}) flamingo.Logger {
	zapFields := make([]zapcore.Field, len(fields))
	area := l.configArea
	traceId := l.traceID
	i := 0

	for key, value := range fields {
		if key == flamingo.LogKeyArea {
			area, _ = value.(string)
		}

		if key == flamingo.LogKeyTraceID {
			traceId, _ = value.(string)
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

	return &SilentLogger{
		Logger:          l.Logger.With(zapFields[:i]...),
		configArea:      area,
		fieldMap:        l.fieldMap,
		logSession:      l.logSession,
		loggingRegistry: l.loggingRegistry,
		traceID:         traceId,
	}
}

// WithField creates a child logger and adds structured context to it.
func (l *SilentLogger) WithField(key flamingo.LogKey, value interface{}) flamingo.Logger {
	traceId := l.traceID

	area := l.configArea
	if key == flamingo.LogKeyArea {
		area, _ = value.(string)
	}

	if key == flamingo.LogKeyTraceID {
		traceId, _ = value.(string)
	}

	if alias, ok := l.fieldMap[string(key)]; ok {
		// skip field
		if alias == "-" {
			return l
		}

		key = flamingo.LogKey(alias)
	}

	return &SilentLogger{
		Logger:          l.Logger.With(zap.Any(string(key), value)),
		configArea:      area,
		fieldMap:        l.fieldMap,
		logSession:      l.logSession,
		loggingRegistry: l.loggingRegistry,
		traceID:         traceId,
	}
}

// Flush is used by buffered loggers and triggers the actual writing. It is a good habit to call Flush before
// letting the process exit. For the top level flamingo.Logger, this is called by the app itself.
func (l *SilentLogger) Flush() {
	_ = l.Logger.Sync()
}
