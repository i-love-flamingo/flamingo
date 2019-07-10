package controllers

import (
	"net/http"

)

type (

	// Version controller
	Version struct {
		versionFile string
	}

)

// Inject Version dependencies
func (h *Version) Inject(config *struct {
	VersionFile     string `inject:"config:healthcheck.versionFile"`
},) {
	if config != nil {
		h.versionFile = config.VersionFile
	}
}


// ServeHTTP responds to Version requests
func (p *Version) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "version.json")
}
