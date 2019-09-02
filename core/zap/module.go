package zap

import (
	"context"
	"fmt"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/zap/application"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Module for logrus logging
	Module struct {
		area               string
		json               bool
		logLevel           string
		coloredOutput      bool
		developmentMode    bool
		samplingEnabled    bool
		samplingInitial    float64
		samplingThereafter float64
		fieldMap           map[string]string
		logSession         bool
	}

	shutdownEventSubscriber struct {
		logger flamingo.Logger
	}
)

var logLevels = map[string]zapcore.Level{
	"Debug":  zap.DebugLevel,
	"Info":   zap.InfoLevel,
	"Warn":   zap.WarnLevel,
	"Error":  zap.ErrorLevel,
	"DPanic": zap.DPanicLevel,
	"Panic":  zap.PanicLevel,
	"Fatal":  zap.FatalLevel,
}

// Inject dependencies
func (m *Module) Inject(config *struct {
	Area               string     `inject:"config:area"`
	JSON               bool       `inject:"config:zap.json,optional"`
	LogLevel           string     `inject:"config:zap.loglevel,optional"`
	ColoredOutput      bool       `inject:"config:zap.colored,optional"`
	DevelopmentMode    bool       `inject:"config:zap.devmode,optional"`
	SamplingEnabled    bool       `inject:"config:zap.sampling.enabled,optional"`
	SamplingInitial    float64    `inject:"config:zap.sampling.initial,optional"`
	SamplingThereafter float64    `inject:"config:zap.sampling.thereafter,optional"`
	FieldMap           config.Map `inject:"config:zap.fieldmap,optional"`
	LogSession         bool       `inject:"config:zap.logsession,optional"`
}) {
	m.area = config.Area
	m.json = config.JSON
	m.logLevel = config.LogLevel
	m.coloredOutput = config.ColoredOutput
	m.developmentMode = config.DevelopmentMode
	m.samplingEnabled = config.SamplingEnabled
	m.samplingInitial = config.SamplingInitial
	m.samplingThereafter = config.SamplingThereafter
	m.logSession = config.LogSession
	if config.FieldMap != nil {
		m.fieldMap = make(map[string]string, len(config.FieldMap))
		for k, v := range config.FieldMap {
			if v, ok := v.(string); ok {
				m.fieldMap[k] = v
			}
		}
	}
}

// Configure the logrus logger as flamingo.Logger (in JSON mode kibana compatible)
func (m *Module) Configure(injector *dingo.Injector) {
	level, ok := logLevels[m.logLevel]
	if !ok {
		// if nothing is configured user ErrorLevel
		level = zap.ErrorLevel
	}

	var samplingConfig *zap.SamplingConfig

	if m.samplingEnabled && m.samplingThereafter > 0 && m.samplingInitial > 0 {
		samplingConfig = &zap.SamplingConfig{
			Initial:    int(m.samplingInitial),
			Thereafter: int(m.samplingThereafter),
		}
	}

	output := "console"
	if m.json {
		output = "json"
	}
	encoder := zapcore.CapitalLevelEncoder
	if m.coloredOutput {
		encoder = zapcore.CapitalColorLevelEncoder
	}
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       m.developmentMode,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          samplingConfig,
		Encoding:          output,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     flamingo.LogKeyMessage,
			LevelKey:       flamingo.LogKeyLevel,
			TimeKey:        flamingo.LogKeyTimestamp,
			NameKey:        "logger",
			CallerKey:      flamingo.LogKeySource,
			StacktraceKey:  flamingo.LogKeyTrace,
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

	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	zapLogger := &Logger{
		Logger:     logger,
		fieldMap:   m.fieldMap,
		logSession: m.logSession,
	}

	zapLogger = zapLogger.WithField(flamingo.LogKeyArea, m.area).(*Logger)

	injector.Bind(new(flamingo.Logger)).ToInstance(zapLogger)
	flamingo.BindEventSubscriber(injector).To(shutdownEventSubscriber{})

	if err := opencensus.View("flamingo/zap/errors", application.ErrorCount, view.Count()); err != nil {
		panic(fmt.Sprintf("failed to register opencensus view: %s", err))
	}
}

// Inject dependencies
func (subscriber *shutdownEventSubscriber) Inject(logger flamingo.Logger) {
	subscriber.logger = logger
}

// Notify handles the incoming event if it is a AppShutdownEvent
func (subscriber *shutdownEventSubscriber) Notify(_ context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		if logger, ok := subscriber.logger.(*Logger); ok {
			logger.Debug("Zap Logger shutdown event")
			logger.Sync()
		}
	}
}

// DefaultConfig for zap log level
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"zap.loglevel":            "Debug",
		"zap.sampling.enabled":    true,
		"zap.sampling.initial":    100,
		"zap.sampling.thereafter": 100,
	}
}
