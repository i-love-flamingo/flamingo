package application

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
)

type (
	// SystemServer provides the internal endpoint
	SystemServer struct {
		handlerProvider domain.HandlerProvider
		logger          flamingo.Logger
		serviceAddress  string
		server          *http.Server
	}
)

// Inject dependencies
func (s *SystemServer) Inject(
	handlerProvider domain.HandlerProvider,
	logger flamingo.Logger,
	config *struct {
		ServiceAddress string `inject:"config:systemendpoint.serviceAddr"`
	},
) {
	s.handlerProvider = handlerProvider
	s.logger = logger
	s.serviceAddress = config.ServiceAddress
}

// Notify handles required actions on startup and shutdown
func (s *SystemServer) Notify(_ context.Context, e flamingo.Event) {
	switch e.(type) {
	case *flamingo.ServerStartEvent:
		s.startup()
	case *flamingo.ServerShutdownEvent:
		s.shutdown()
	}
}

func (s *SystemServer) startup() {
	s.logger.Info("systemendpoint: startup at ", s.serviceAddress)
	serveMux := http.NewServeMux()
	for route, handler := range s.handlerProvider() {
		if handler != nil {
			s.logger.Debug("systemendpoint: register route ", route)
			serveMux.Handle(route, handler)
		}
	}
	s.server = &http.Server{Addr: s.serviceAddress, Handler: serveMux}
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}

func (s *SystemServer) shutdown() {
	s.logger.Info("systemendpoint: shutdown at ", s.serviceAddress)
	if s.server != nil {
		_ = s.server.Shutdown(context.Background())
	}
}
