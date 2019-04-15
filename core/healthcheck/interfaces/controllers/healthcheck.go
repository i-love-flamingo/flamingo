package controllers

import (
	"encoding/json"
	"net/http"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
)

type (
	statusProvider func() map[string]healthcheck.Status

	// Healthcheck controller
	Healthcheck struct {
		statusProvider statusProvider
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
func (h *Healthcheck) Inject(provider statusProvider) {
	h.statusProvider = provider
}

// ServeHTTP responds to healthcheck requests
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

// ServeHTTP responds to Ping requests
func (p *Ping) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	handleErr(err, w)
}

// TryServeHTTP implementation to be used in prefixrouter & co
func (p *Ping) TryServeHTTP(rw http.ResponseWriter, req *http.Request) (bool, error) {
	if req.URL.Path != "/health/ping" {
		return true, nil
	}

	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("OK"))
	return false, err
}

func handleErr(err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}
