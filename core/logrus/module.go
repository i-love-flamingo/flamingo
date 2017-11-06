package logrus

import (
	"github.com/sirupsen/logrus"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/flamingo"
	"os"
)

type (
	// Module for logrus logging
	Module struct {
		Area string `inject:"config:area"`
		JSON bool   `inject:"config:logrus.json,optional"`
	}

	logrusEntry struct {
		*logrus.Entry
	}
)

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
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(&logrusEntry{l.WithField("area", m.Area)})
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
