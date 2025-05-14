package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/opencensus"
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

const (
	statusMeasureChecksName   = "flamingo/status/checks"
	statusMeasureFailuresName = "flamingo/status/failures"
)

var (
	// healthcheckStatusFailureMeasure counts failures of status checks
	healthcheckStatusChecksMeasure  = stats.Int64(statusMeasureChecksName, "Count of status checks.", stats.UnitDimensionless)
	healthcheckStatusFailureMeasure = stats.Int64(statusMeasureFailuresName, "Count of status check failures.", stats.UnitDimensionless)

	serviceName, _ = tag.NewKey("service_name")
)

func init() {
	if err := opencensus.View(statusMeasureFailuresName, healthcheckStatusFailureMeasure, view.Count(), serviceName); err != nil {
		panic(err)
	}

	if err := opencensus.View(statusMeasureChecksName, healthcheckStatusChecksMeasure, view.Count(), serviceName); err != nil {
		panic(err)
	}
}

// Inject Healthcheck dependencies
func (h *Healthcheck) Inject(provider statusProvider) {
	h.statusProvider = provider
}

// ServeHTTP responds to healthcheck requests
func (h *Healthcheck) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var resp response
	var allAlive = true

	for name, status := range h.statusProvider() {
		increaseMeasureIfMeasuredStatus(req.Context(), status, healthcheckStatusChecksMeasure)

		alive, details := status.Status()
		if !alive {
			increaseMeasureIfMeasuredStatus(req.Context(), status, healthcheckStatusFailureMeasure)

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

func increaseMeasureIfMeasuredStatus(ctx context.Context, status healthcheck.Status, measure *stats.Int64Measure) {
	if measuredStatus, ok := status.(healthcheck.MeasuredStatus); ok {
		recordContext, _ := tag.New(ctx, tag.Upsert(serviceName, measuredStatus.ServiceName()))

		stats.Record(recordContext, measure.M(1))
	}
}
