package controllers

import (
	"context"

	"net/http"
	"strings"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	StatusProvider func() map[string]healthcheck.Status

	Healthcheck struct {
		responder      *web.Responder
		statusProvider StatusProvider
	}

	response struct {
		Services []service `json:"services,omitempty"`
	}

	service struct {
		Name    string `json:"name"`
		Alive   bool   `json:"alive"`
		Details string `json:"details"`
	}
)

// Inject Healthcheck dependencies
func (controller *Healthcheck) Inject(responder *web.Responder, provider StatusProvider) {
	controller.responder = responder
	controller.statusProvider = provider
}

func (controller *Healthcheck) Healthcheck(context.Context, *web.Request) web.Response {
	var resp response

	for name, status := range controller.statusProvider() {
		alive, details := status.Status()

		resp.Services = append(resp.Services, service{
			Name:    name,
			Alive:   alive,
			Details: details,
		})
	}

	return controller.responder.Data(resp)
}

func (controller *Healthcheck) Ping(context.Context, *web.Request) web.Response {
	return controller.responder.HTTP(uint(http.StatusOK), strings.NewReader("OK"))
}
