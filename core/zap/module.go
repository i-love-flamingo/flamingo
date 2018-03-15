package zap

import (
	"go.aoe.com/flamingo/core/zap/domain"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Module for logrus logging
	Module struct {
		Area               string  `inject:"config:area"`
		JSON               bool    `inject:"config:zap.json,optional"`
		LogLevel           string  `inject:"config:zap.loglevel,optional"`
		DevelopmentMode    bool    `inject:"config:zap.devmode,optional"`
		SamplingEnabled    bool    `inject:"config:zap.sampling.enabled,optional"`
		SamplingInitial    float32 `inject:"config:zap.sampling.initial,optional"`
		SamplingThereafter float32 `inject:"config:zap.sampling.thereafter,optional"`
	}

	// ShutdownEventSubscriber handles the logger sync on flamingo shutdown
	ShutdownEventSubscriber struct {
		Logger flamingo.Logger `inject:""`
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

// Configure the logrus logger as flamingo.Logger (in JSON mode kibana compatible)
func (m *Module) Configure(injector *dingo.Injector) {
	level, ok := logLevels[m.LogLevel]
	if !ok {
		// if nothing is configured user ErrorLevel
		level = zap.ErrorLevel
	}

	var samplingConfig *zap.SamplingConfig

	if m.SamplingEnabled {
		samplingConfig = &zap.SamplingConfig{
			Initial:    int(m.SamplingInitial),
			Thereafter: int(m.SamplingThereafter),
		}
	}

	output := "console"
	if m.JSON {
		output = "json"
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       m.DevelopmentMode,
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
			EncodeLevel:    zapcore.CapitalLevelEncoder,
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
	logger = logger.With(zap.String(flamingo.LogKeyArea, m.Area))

	zapLogger := &domain.Logger{
		Logger: logger,
	}

	injector.Bind((*flamingo.Logger)(nil)).ToInstance(zapLogger)
	injector.BindMulti((*event.Subscriber)(nil)).To(ShutdownEventSubscriber{})
}

// Notify handles the incoming event if it is a AppShutdownEvent
func (subscriber *ShutdownEventSubscriber) Notify(event event.Event) {
	switch event.(type) {
	case *flamingo.AppShutdownEvent:
		if logger, ok := subscriber.Logger.(*domain.Logger); ok {
			logger.Debug("Zap Logger shutdown event")
			logger.Sync()
		}
	}
}
