package silentzap

import (
	"context"
	"fmt"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// Module for Silent ZAP logging
	Module struct{}

	shutdownEventSubscriber struct {
		logger flamingo.Logger
	}
)

// Configure the ZAP logger as flamingo.Logger (in JSON mode kibana compatible)
func (m *Module) Configure(injector *dingo.Injector) {
	registry := new(LoggingContextRegistry)

	injector.Bind(registry).AsEagerSingleton()

	injector.Bind(new(flamingo.Logger)).ToProvider(getSilentLogger)

	flamingo.BindEventSubscriber(injector).To(shutdownEventSubscriber{})
	flamingo.BindEventSubscriber(injector).To(registry)
}

// Inject dependencies
func (subscriber *shutdownEventSubscriber) Inject(logger flamingo.Logger) {
	subscriber.logger = logger
}

// Notify handles the incoming event if it is a AppShutdownEvent
func (subscriber *shutdownEventSubscriber) Notify(_ context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		if logger, ok := subscriber.logger.(*SilentLogger); ok {
			//nolint:contextcheck //no need for context here
			logger.Debug("Silent Zap Logger shutdown event")
			_ = logger.Sync()
		}
	}
}

// CueConfig Schema
func (m *Module) CueConfig() string {
	// language=cue
	return fmt.Sprintf(`
core zap: {
	loglevel: %s
	sampling: {
		enabled: bool | *true
		initial: int | *100 
		thereafter: int | *100
	}
}
`, allowedLevels())
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
