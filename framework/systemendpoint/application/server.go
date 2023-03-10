package application

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"
	"sync"

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
		mu              sync.Mutex
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

// Start the systemendpoint in a separate go routine
func (s *SystemServer) Start() {
	serveMux := http.NewServeMux()
	for route, handler := range s.handlerProvider() {
		if handler != nil {
			s.logger.Debug("systemendpoint: register route ", route)
			serveMux.Handle(route, handler)
		}
	}

	serveMux.HandleFunc("/version", func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(writer, "version: %s\n", flamingo.AppVersion())
		_, _ = fmt.Fprintf(writer, "go: %s\n", runtime.Version())
		if info, ok := debug.ReadBuildInfo(); ok {
			_, _ = fmt.Fprintf(writer, "path: %s\n", info.Path)
			for _, module := range info.Deps {
				if module.Path == "flamingo.me/flamingo/v3" {
					_, _ = fmt.Fprintf(writer, "flamingo: %s\n", module.Version)
				}
			}
		}
	})

	listener, err := net.Listen("tcp", s.serviceAddress)
	if err != nil {
		s.logger.Fatal(err)
	}

	s.mu.Lock()
	s.server = &http.Server{Handler: serveMux}
	s.mu.Unlock()
	go func() {
		s.logger.Info("Starting HTTP Server (systemendpoint) at ", listener.Addr())
		err := s.server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (s *SystemServer) shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.server != nil {
		s.logger.Info("systemendpoint: shutdown")
		_ = s.server.Shutdown(context.Background())
		s.server = nil
	}
}
