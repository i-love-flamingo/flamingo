package logrus

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	// Module for logrus logging
	Module struct {
		Area string `inject:"config:area"`
		JSON bool   `inject:"config:logrus.json,optional"`
	}

	LogrusEntry struct {
		*logrus.Entry
	}

	ContextHook struct{}
)

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(6, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["file"] = path.Base(file)
			entry.Data["func"] = path.Base(name)
			entry.Data["line"] = line
			break
		}
	}
	return nil
}

func (m *Module) Configure(injector *dingo.Injector) {
	var l *logrus.Logger
	if m.JSON {
		l = &logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		}
	} else {
		l = logrus.New()
	}
	l.Hooks.Add(ContextHook{})
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(&LogrusEntry{l.WithField("area", m.Area)})
}

func (e *LogrusEntry) WithField(key string, value interface{}) flamingo.Logger {
	return &LogrusEntry{e.Entry.WithField(key, value)}
}

func (e *LogrusEntry) WithFields(fields map[string]interface{}) flamingo.Logger {
	return &LogrusEntry{e.Entry.WithFields(fields)}
}

func (e *LogrusEntry) WithError(err error) flamingo.Logger {
	return &LogrusEntry{e.Entry.WithError(err)}
}
