package process_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthcheck_FailingStatus(t *testing.T) {
	t.Parallel()

	_, addrs := serve(t, testappMain, nil, "--flamingo-config", "ops.health.fail: true")

	got := get(t, addrs.System, "/status/healthcheck")
	assert.Equal(t, 500, got.status)

	var raw map[string]any

	require.NoError(t, json.Unmarshal([]byte(got.body), &raw), got.body)

	services, _ := raw["services"].([]any)
	entries := map[string]map[string]any{}

	for _, item := range services {
		service, _ := item.(map[string]any)
		name, _ := service["name"].(string)
		entries[name] = service
	}

	// Decoding into maps keeps the check on the letter case of the keys.
	assert.Equal(t, map[string]any{"name": "ops", "alive": false, "details": "ops down"}, entries["ops"])

	if assert.Contains(t, entries, "session") {
		assert.Equal(t, true, entries["session"]["alive"])
	}
}
