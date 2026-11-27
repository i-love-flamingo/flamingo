package metrics_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats/view"

	// The blank imports load the packages that register the views.
	_ "flamingo.me/flamingo/v3/core/healthcheck/interfaces/controllers"
	_ "flamingo.me/flamingo/v3/core/oauth/interfaces"
	_ "flamingo.me/flamingo/v3/core/silentzap"
	_ "flamingo.me/flamingo/v3/core/zap"
	_ "flamingo.me/flamingo/v3/framework/prefixrouter"
	_ "flamingo.me/flamingo/v3/framework/web"
	_ "flamingo.me/flamingo/v3/framework/web/filter"
)

// Dashboards and alerts select metrics by view name and group them by tag key.
func TestViews_NamesAndTagKeys(t *testing.T) {
	t.Parallel()

	want := map[string][]string{
		"flamingo/router/controller":     {"area", "controller", "method"},
		"flamingo/response/bytes_count":  {"area", "status_code"},
		"flamingo/response/bytes_sum":    {"area", "status_code"},
		"flamingo/prefixrouter/requests": {"area"},
		"flamingo/zap/logs":              {"area", "level"},
		"flamingo/silentlogger/logs":     {"area", "level"},
		"flamingo/oauth/login_failed":    {"area"},
		"flamingo/oauth/login_succeeded": {"area"},
		"flamingo/status/checks/total":   {"area", "name"},
		"flamingo/status/failures/total": {"area", "name"},
	}

	for name, keys := range want {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			registered := view.Find(name)
			require.NotNil(t, registered, "view is not registered")

			got := make([]string, 0, len(registered.TagKeys))
			for _, key := range registered.TagKeys {
				got = append(got, key.Name())
			}

			slices.Sort(got)

			assert.Equal(t, keys, got)
		})
	}
}
