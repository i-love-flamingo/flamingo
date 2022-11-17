package zap

import (
	"context"
	"fmt"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Module for zap logging
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
		callerEncoder      zapcore.CallerEncoder
	}

	shutdownEventSubscriber struct {
		logger flamingo.Logger
	}
)

const (
	ZapCallerEncoderShort = "short"
	ZapCallerEncoderSmart = "smart"
	ZapCallerEncoderFull  = "full"
)

var logLevels = map[string]zapcore.Level{
	"Trace":  zap.DebugLevel - 1, // does not exist in zap by default
	"Debug":  zap.DebugLevel,
	"Info":   zap.InfoLevel,
	"Warn":   zap.WarnLevel,
	"Error":  zap.ErrorLevel,
	"DPanic": zap.DPanicLevel,
	"Panic":  zap.PanicLevel,
	"Fatal":  zap.FatalLevel,
}

var callerEncoders = map[string]zapcore.CallerEncoder{
	ZapCallerEncoderSmart: smartCallerEncoder,
	ZapCallerEncoderFull:  zapcore.FullCallerEncoder,
	ZapCallerEncoderShort: zapcore.ShortCallerEncoder,
}

// Inject dependencies
func (m *Module) Inject(config *struct {
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
	CallerEncoder      string     `inject:"config:core.zap.encoding.caller,optional"`
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
	m.callerEncoder = callerEncoders[ZapCallerEncoderShort]

	if encoder, ok := callerEncoders[config.CallerEncoder]; ok {
		m.callerEncoder = encoder
	}

	if config.FieldMap != nil {
		m.fieldMap = make(map[string]string, len(config.FieldMap))

		for k, v := range config.FieldMap {
			if v, ok := v.(string); ok {
				m.fieldMap[k] = v
			}
		}
	}
}

// Configure the zap logger as flamingo.Logger
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(flamingo.Logger)).ToInstance(m.createLoggerInstance())
	flamingo.BindEventSubscriber(injector).To(shutdownEventSubscriber{})
}

func (m *Module) createLoggerInstance() *Logger {
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

	// Capital encoder with trace addition
	encoder := capitalLevelEncoder
	if m.coloredOutput {
		// Capital color encoder with trace addition
		encoder = capitalColorLevelEncoder
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       m.developmentMode,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          samplingConfig,
		Encoding:          output,
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
			EncodeCaller:   m.callerEncoder,
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

	zapLogger := NewLogger(
		logger,
		WithFieldMap(m.fieldMap),
		WithLogSession(m.logSession),
		WithArea(m.area),
	)

	zapLogger = zapLogger.WithField(flamingo.LogKeyArea, m.area).(*Logger)

	return zapLogger
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
			_ = logger.Sync()
		}
	}
}

// CueConfig Schema
func (m *Module) CueConfig() string {
	// language=cue
	return fmt.Sprintf(`
core: zap: {
	json: bool | *false
	colored: bool | *false
	devmode: bool | *false
	logsession: bool | *false
	fieldmap: {
		[string]: string
	}
	
	loglevel: %s
	sampling: {
		enabled: bool | *true
		initial: int | *100 
		thereafter: int | *100
	}
	
	encoding: {
		caller: *"%s" | "%s" | "%s"
	}
}
`, allowedLevels, ZapCallerEncoderShort, ZapCallerEncoderSmart, ZapCallerEncoderFull)
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"zap.loglevel":            "core.zap.loglevel",
		"zap.sampling.enabled":    "core.zap.sampling.enabled",
		"zap.sampling.initial":    "core.zap.sampling.initial",
		"zap.sampling.thereafter": "core.zap.sampling.thereafter",
	}
}
