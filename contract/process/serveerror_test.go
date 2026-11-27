package process_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const superfluous = "superfluous response.WriteHeader"

// An error that net/http logs, here a second WriteHeader, becomes an ERROR line of the JSON log.
// This is flamingo's own serve command, without prefixrouter.
func TestServeErrors_ServeCommandUsesLogger(t *testing.T) {
	t.Parallel()

	p, addrs := serve(t, plainserve, nil)

	require.Equal(t, 200, get(t, addrs.HTTP, "/ops/twice").status)

	require.Len(t, linesWith(t, p, superfluous), 1)

	var found []map[string]any

	for _, entry := range p.Logs() {
		if message, _ := entry["message"].(string); strings.Contains(message, superfluous) {
			found = append(found, entry)
		}
	}

	require.Len(t, found, 1, "JSON lines with the error")
	assert.Equal(t, "ERROR", found[0]["level"])
}
