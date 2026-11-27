package process_test

import (
	"regexp"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/contract/internal/harness"
)

const eventPrefix = "contract event *flamingo."

var mainServerLine = regexp.MustCompile(`^Starting HTTP Server \(Prefixrouter\) at (\S+)$`)

// waitMainServer waits for the log line of the main server, and not for the system endpoint, which starts later.
func waitMainServer(tb testing.TB, p *harness.Proc) string {
	tb.Helper()

	var addr string

	require.Eventually(tb, func() bool {
		for _, entry := range p.Logs() {
			message, _ := entry["message"].(string)

			if match := mainServerLine.FindStringSubmatch(message); match != nil {
				addr = match[1]

				return true
			}
		}

		return false
	}, 30*time.Second, 10*time.Millisecond, "no line of the main server, lines: %v", p.Lines())

	return addr
}

// countEvents gives how often the log has the line of a lifecycle event.
func countEvents(p *harness.Proc, event string) int {
	var count int

	for _, entry := range p.Logs() {
		if entry["message"] == eventPrefix+event {
			count++
		}
	}

	return count
}

func TestLifecycle_EventsAndSIGTERM(t *testing.T) {
	t.Parallel()

	p := harness.Process(t, harness.Options{
		Main: testappMain,
		Dir:  "../internal/testapp",
		Args: []string{"--flamingo-config", "ops.events.delayMS: 200", "serve", "--addr", "127.0.0.1:0"},
	})

	// The subscriber of the app is slow. The server must still not answer before the start events are handled.
	started := get(t, waitMainServer(t, p), "/en/ops/started")
	require.Equal(t, 200, started.status)
	assert.JSONEq(t, `{"startup": true, "serverStart": true}`, started.body)

	// the event lines can come after the response, so the test polls the log
	require.Eventually(t, func() bool {
		return countEvents(p, "StartupEvent") == 1 && countEvents(p, "ServerStartEvent") == 1
	}, 30*time.Second, 10*time.Millisecond, "events at the start, lines: %v", p.Lines())

	assert.Zero(t, countEvents(p, "ServerShutdownEvent"))
	assert.Zero(t, countEvents(p, "ShutdownEvent"))

	p.Signal(t, syscall.SIGTERM)

	assert.Equal(t, 0, p.Wait(t))
	assert.Equal(t, 1, countEvents(p, "ServerShutdownEvent"))
	assert.Equal(t, 1, countEvents(p, "ShutdownEvent"))
	assert.Equal(t, 1, countEvents(p, "StartupEvent"))
	assert.Equal(t, 1, countEvents(p, "ServerStartEvent"))
}

func TestLifecycle_CommandReadsContextConfig(t *testing.T) {
	t.Parallel()

	p := harness.Process(t, harness.Options{
		Main: testappMain,
		Dir:  "../internal/testapp",
		Env:  map[string]string{"CONTEXT": "ci"},
		Args: []string{"ops-greet"},
	})

	assert.Equal(t, 0, p.Wait(t))
	assert.Contains(t, p.Lines(), "ci")
}
