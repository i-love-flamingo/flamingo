package logrus

import (
	"fmt"
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
		Area     string `inject:"config:area"`
		JSON     bool   `inject:"config:logrus.json,optional"`
		LogLevel string `inject:"config:logrus.loglevel,optional"`
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

// Levels returns all available logrus log levels
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
			entry.Data["source"] = fmt.Sprintf("File: %v Func: %v  Line: %v", lastPathAndFileName(file), path.Base(name), line)
			entry.Data["fileName"] = path.Base(file)
			break
		}
	}
	return nil
}

//lastPathAndFileName returns the filename and the last two folders of the Path
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
	var l *logrus.Logger
	if m.JSON {
		l = &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "@timestamp",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
				},
			},
			Hooks: make(logrus.LevelHooks),
		}
	} else {
		l = &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			},
			Hooks: make(logrus.LevelHooks),
		}
	}
	l.Level = logrus.ErrorLevel
	if m.LogLevel == "Info" {
		l.Level = logrus.InfoLevel
	} else if m.LogLevel == "Debug" {
		l.Level = logrus.DebugLevel
	}
	l.Hooks.Add(ContextHook{area: m.Area})
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(&LogrusLogger{l})
}

// WithField adds a single field to the Entry.
func (e *LogrusEntry) WithField(key string, value interface{}) flamingo.Logger {
	return &LogrusEntry{e.Entry.WithField(key, value)}
}

// WithFields adds a map of fields to the Entry.
func (e *LogrusEntry) WithFields(fields map[string]interface{}) flamingo.Logger {
	return &LogrusEntry{e.Entry.WithFields(fields)}
}

// WithError adds an error as single field (using the key defined in ErrorKey) to the Entry.
func (e *LogrusEntry) WithError(err error) flamingo.Logger {
	return &LogrusEntry{e.Entry.WithError(err)}
}

// WithField adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (e *LogrusLogger) WithField(key string, value interface{}) flamingo.Logger {
	return &LogrusEntry{e.Logger.WithField(key, value)}
}

// WithFields adds a struct of fields to the log entry. All it does is call `WithField` for each `Field`.
func (e *LogrusLogger) WithFields(fields map[string]interface{}) flamingo.Logger {
	return &LogrusEntry{e.Logger.WithFields(fields)}
}

// WithError add an error as single field to the log entry.  All it does is call `WithError` for the given `error`.
func (e *LogrusLogger) WithError(err error) flamingo.Logger {
	return &LogrusEntry{e.Logger.WithError(err)}
}
