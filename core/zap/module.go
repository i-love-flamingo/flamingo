package logrus

import (
	"path"
	"strings"

	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Module for logrus logging
	Module struct {
		Area     string `inject:"config:area"`
		JSON     bool   `inject:"config:zap.json,optional"`
		LogLevel string `inject:"config:zap.loglevel,optional"`
	}
)

//lastPathAndFileName returns the filename and the last two folders of the Path
//entry.Data["source"] = fmt.Sprintf("File: %v Func: %v  Line: %v", lastPathAndFileName(file), path.Base(name), line)
//entry.Data["fileName"] = path.Base(file)
func lastPathAndFileName(completePath string) string {
	path, fileName := path.Split(completePath)
	dirNames := strings.Split(strings.Trim(path, "/"), "/")
	if len(dirNames) > 2 {
		return dirNames[len(dirNames)-2] + "/" + dirNames[len(dirNames)-1] + "/" + fileName
	}
	return fileName
}

// Configure the logrus logger as flamingo.Logger (in JSON mode kibana compatible)
func (m *Module) Configure(injector *dingo.Injector) {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json", //"console"
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "source",
			MessageKey:     "message",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
		OutputPaths:      []string{"stderr"}, //stdout
		ErrorOutputPaths: []string{"stderr"}, //stdout
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zapLogger := &Logger{
		Logger: logger,
	}

	injector.Bind((*flamingo.Logger)(nil)).ToInstance(zapLogger)
}
