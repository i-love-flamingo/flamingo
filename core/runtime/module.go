package runtime

import (
	"fmt"
	"sync"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"go.uber.org/automaxprocs/maxprocs"
)

type (
	// Module for core.runtime
	Module struct {
		logger flamingo.Logger
	}
)

var once = sync.Once{}

// Inject Module dependencies
func (m *Module) Inject(logger flamingo.Logger) *Module {
	m.logger = logger.WithField(flamingo.LogKeyModule, "core.runtime").WithField(flamingo.LogKeyCategory, "module")
	return m
}

// Configure runtime dependency injection
func (m *Module) Configure(injector *dingo.Injector) {
	once.Do(func() {
		_, err := maxprocs.Set(maxprocs.Logger(
			func(logMessage string, args ...interface{}) {
				m.logger.Info(fmt.Sprintf(logMessage, args...))
			},
		))

		if err != nil {
			m.logger.Error(fmt.Sprintf("Failed to set maxprocs: %v", err))
		}
	})
}
