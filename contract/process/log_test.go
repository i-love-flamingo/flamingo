package process_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Log searches of applications depend on these keys.
func TestLog_JSONLinesOnStderr(t *testing.T) {
	t.Parallel()

	p, _ := serve(t, testappMain, nil)

	for _, line := range p.Lines() {
		assert.True(t, json.Valid([]byte(line)), "not a JSON line: %s", line)
	}

	logs := p.Logs()
	require.NotEmpty(t, logs)

	var prefixrouterLines, runtimeLines int

	for _, entry := range logs {
		assert.Contains(t, []any{"INFO", "WARN", "ERROR"}, entry["level"], "level in %v", entry)

		for _, key := range []string{"message", "@timestamp", "source", "area"} {
			assert.NotEmpty(t, entry[key], "key %q in %v", key, entry)
		}

		source, _ := entry["source"].(string)

		if strings.HasPrefix(source, "prefixrouter/") {
			prefixrouterLines++

			assert.Equal(t, "prefixrouter", entry["category"], "category in %v", entry)
		}

		if strings.HasPrefix(source, "runtime/") {
			runtimeLines++

			assert.Equal(t, "core.runtime", entry["module"], "module in %v", entry)
		}
	}

	// the checks above need lines to look at
	assert.Positive(t, prefixrouterLines, "lines from the prefixrouter")
	assert.Positive(t, runtimeLines, "lines from the runtime module")
}
