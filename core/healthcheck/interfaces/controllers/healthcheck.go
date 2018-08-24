package controllers

import (
	"context"

	"flamingo.me/flamingo/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	StatusProvider func() map[string]healthcheck.Status

	Healthcheck struct {
		responder.JSONAware
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
func (controller *Healthcheck) Inject(aware responder.JSONAware, provider StatusProvider) {
	controller.JSONAware = aware
	controller.statusProvider = provider
}

func (controller *Healthcheck) Get(context.Context, *web.Request) web.Response {
	var resp response

	for name, status := range controller.statusProvider() {
		alive, details := status.Status()

		resp.Services = append(resp.Services, service{
			Name:    name,
			Alive:   alive,
			Details: details,
		})
	}

	return controller.JSON(resp)
}
