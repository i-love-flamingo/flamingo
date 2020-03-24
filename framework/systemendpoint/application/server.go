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
		ServiceAddress string `inject:"config:flamingo.systemendpoint.serviceAddr"`
	},
) {
	s.handlerProvider = handlerProvider
	s.logger = logger
	s.serviceAddress = config.ServiceAddress
}

// Notify handles required actions on Start and shutdown
func (s *SystemServer) Notify(_ context.Context, e flamingo.Event) {
	switch e.(type) {
	case *flamingo.ServerStartEvent:
		s.Start()
	case *flamingo.ServerShutdownEvent:
		s.shutdown()
	case *flamingo.ShutdownEvent:
		s.shutdown()
	}
}

//Start - starts the systemendpoint in a sperate go routine
func (s *SystemServer) Start() {
	s.logger.Info("systemendpoint: Start at ", s.serviceAddress)
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
		if err != nil && err != http.ErrServerClosed {
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
