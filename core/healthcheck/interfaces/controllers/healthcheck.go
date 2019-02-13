package controllers

import (
	"encoding/json"
	"net/http"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
)

type (
	// StatusProvider provides all bound healthchecks
	StatusProvider func() map[string]healthcheck.Status

	// Healthcheck controller
	Healthcheck struct {
		statusProvider StatusProvider
	}

	// Ping controller
	Ping struct{}

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
func (h *Healthcheck) Inject(provider StatusProvider) {
	h.statusProvider = provider
}

func (h *Healthcheck) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	var resp response
	var allAlive = true

	for name, status := range h.statusProvider() {
		alive, details := status.Status()
		if !alive {
			allAlive = false
		}

		resp.Services = append(resp.Services, service{
			Name:    name,
			Alive:   alive,
			Details: details,
		})
	}

	var status = http.StatusOK
	if !allAlive {
		status = http.StatusInternalServerError
	}

	respBody, err := json.Marshal(resp)
	handleErr(err, w)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(respBody)
	handleErr(err, w)
}

func (p *Ping) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	handleErr(err, w)
}

func handleErr(err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}
