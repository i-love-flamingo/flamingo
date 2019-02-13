package zap

import (
	"context"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
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
		samplingInitial    float32
		samplingThereafter float32
	}

	// ShutdownEventSubscriber handles the logger sync on flamingo shutdown
	ShutdownEventSubscriber struct {
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
	Area               string  `inject:"config:area"`
	JSON               bool    `inject:"config:zap.json,optional"`
	LogLevel           string  `inject:"config:zap.loglevel,optional"`
	ColoredOutput      bool    `inject:"config:zap.colored,optional"`
	DevelopmentMode    bool    `inject:"config:zap.devmode,optional"`
	SamplingEnabled    bool    `inject:"config:zap.sampling.enabled,optional"`
	SamplingInitial    float32 `inject:"config:zap.sampling.initial,optional"`
	SamplingThereafter float32 `inject:"config:zap.sampling.thereafter,optional"`
}) {
	m.area = config.Area
	m.json = config.JSON
	m.logLevel = config.LogLevel
	m.coloredOutput = config.ColoredOutput
	m.developmentMode = config.DevelopmentMode
	m.samplingEnabled = config.SamplingEnabled
	m.samplingInitial = config.SamplingInitial
	m.samplingThereafter = config.SamplingThereafter
}

// Configure the logrus logger as flamingo.Logger (in JSON mode kibana compatible)
func (m *Module) Configure(injector *dingo.Injector) {
	level, ok := logLevels[m.logLevel]
	if !ok {
		// if nothing is configured user ErrorLevel
		level = zap.ErrorLevel
	}

	var samplingConfig *zap.SamplingConfig

	if m.samplingEnabled {
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
			EncodeCaller:   zapcore.FullCallerEncoder,
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
	logger = logger.With(zap.String(flamingo.LogKeyArea, m.area))

	zapLogger := &Logger{
		Logger: logger,
	}

	injector.Bind(new(flamingo.Logger)).ToInstance(zapLogger)
	flamingo.BindEventSubscriber(injector).To(ShutdownEventSubscriber{})
}

// Inject dependencies
func (subscriber *ShutdownEventSubscriber) Inject(logger flamingo.Logger) {
	subscriber.logger = logger
}

// Notify handles the incoming event if it is a AppShutdownEvent
func (subscriber *ShutdownEventSubscriber) Notify(_ context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		if logger, ok := subscriber.logger.(*Logger); ok {
			logger.Debug("Zap Logger shutdown event")
			_ = logger.Sync()
		}
	}
}

// DefaultConfig for zap log level
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"zap.loglevel": "Debug",
	}
}
