package logrus

import (
	"os"
	"path"
	"runtime"
	"strings"

	"sync"

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

	LogrusLogger struct {
		*logrus.Logger
	}

	ContextHook struct {
		area string
	}
)

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var lock = new(sync.Mutex)

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	lock.Lock()
	defer lock.Unlock()

	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(6, pc)

	entry.Data["area"] = hook.area

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
	l.Hooks.Add(ContextHook{area: m.Area})
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(&LogrusLogger{l})
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

func (e *LogrusLogger) WithField(key string, value interface{}) flamingo.Logger {
	return &LogrusEntry{e.Logger.WithField(key, value)}
}

func (e *LogrusLogger) WithFields(fields map[string]interface{}) flamingo.Logger {
	return &LogrusEntry{e.Logger.WithFields(fields)}
}

func (e *LogrusLogger) WithError(err error) flamingo.Logger {
	return &LogrusEntry{e.Logger.WithError(err)}
}
