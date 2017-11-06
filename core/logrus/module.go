package logrus

import (
	"github.com/sirupsen/logrus"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/flamingo"
	"os"
)

// Module for logrus logging
type (
	Module struct {
		Area string `inject:"config:area"`
		JSON bool   `inject:"config:logrus.json,optional"`
	}

	logrusLogger struct {
		*logrus.Logger
	}

	logrusEntry struct {
		*logrus.Entry
	}
)

func (m *Module) Configure(injector *dingo.Injector) {
	var l flamingo.Logger
	if m.JSON {
		l = &logrusLogger{
			Logger: &logrus.Logger{
				Out:       os.Stderr,
				Formatter: new(logrus.JSONFormatter),
				Hooks:     make(logrus.LevelHooks),
				Level:     logrus.InfoLevel,
			},
		}
	} else {
		l = &logrusLogger{
			Logger: logrus.New(),
		}
	}
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(l.WithField("area", m.Area))
}

func (l *logrusLogger) WithField(key string, value interface{}) flamingo.Logger {
	return &logrusEntry{l.Logger.WithField(key, value)}
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) flamingo.Logger {
	return &logrusEntry{l.Logger.WithFields(fields)}
}

func (l *logrusLogger) WithError(err error) flamingo.Logger {
	return &logrusEntry{l.Logger.WithError(err)}
}

func (e *logrusEntry) WithField(key string, value interface{}) flamingo.Logger {
	return &logrusEntry{e.Entry.WithField(key, value)}
}

func (e *logrusEntry) WithFields(fields map[string]interface{}) flamingo.Logger {
	return &logrusEntry{e.Entry.WithFields(fields)}
}

func (e *logrusEntry) WithError(err error) flamingo.Logger {
	return &logrusEntry{e.Entry.WithError(err)}
}
