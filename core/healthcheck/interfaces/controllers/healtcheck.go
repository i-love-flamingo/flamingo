package controllers

import (
	"flamingo.me/flamingo/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	StatusProvider func() map[string]healthcheck.Status

	Healthcheck struct {
		responder.JSONAware `inject:""`
		StatusProvider      StatusProvider `inject:""`
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

func (controller *Healthcheck) Get(ctx web.Context) web.Response {
	var resp response

	for name, status := range controller.StatusProvider() {
		alive, details := status.Status()

		resp.Services = append(resp.Services, service{
			Name:    name,
			Alive:   alive,
			Details: details,
		})
	}

	return controller.JSON(resp)
}
