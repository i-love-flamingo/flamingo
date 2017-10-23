package logrus

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/flamingo"
	"github.com/sirupsen/logrus"
	"os"
)

// Module for logrus logging
type Module struct {
	Area string `inject:"config:area"`
	JSON bool `inject:"config:logrus.json,optional"`
}

func (m *Module) Configure(injector *dingo.Injector) {
	if m.JSON {
		l := &logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		}
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(l.WithField("area", m.Area))
	} else {
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(logrus.New().WithField("area", m.Area))
	}
}
