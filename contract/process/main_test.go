package process_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp"
)

const (
	testappMain = "testapp"
	plainserve  = "plainserve"
)

var client = &http.Client{Timeout: 30 * time.Second}

func TestMain(m *testing.M) {
	harness.Main(m, map[string]func(){
		testappMain: testapp.Main,
		plainserve:  testapp.PlainServe,
	})
}

// containsLine reports whether one of the lines contains the text.
func containsLine(lines []string, substr string) bool {
	for _, line := range lines {
		if strings.Contains(line, substr) {
			return true
		}
	}

	return false
}

// linesWith waits until a line of stderr contains the text and gives all lines that contain it. The reader of
// stderr can be late, so a test must not read the lines only once right after a request.
func linesWith(tb testing.TB, p *harness.Proc, substr string) []string {
	tb.Helper()

	require.Eventually(tb, func() bool { return containsLine(p.Lines(), substr) }, 30*time.Second, 10*time.Millisecond,
		"no line with %q", substr)

	var found []string

	for _, line := range p.Lines() {
		if strings.Contains(line, substr) {
			found = append(found, line)
		}
	}

	return found
}

// serve starts the app with the command "serve" on a free port and waits for the servers.
func serve(tb testing.TB, app string, env map[string]string, flags ...string) (*harness.Proc, harness.Addrs) {
	tb.Helper()

	args := append(append([]string{}, flags...), "serve", "--addr", "127.0.0.1:0")

	proc := harness.Process(tb, harness.Options{Main: app, Dir: "../internal/testapp", Env: env, Args: args})

	return proc, proc.WaitReady(tb)
}

// reply is the status, the body and the cookies of a response. The body is already read and closed.
type reply struct {
	status  int
	body    string
	cookies []*http.Cookie
}

// get sends a GET request to the address.
func get(tb testing.TB, addr, path string, cookies ...*http.Cookie) reply {
	tb.Helper()

	req, err := http.NewRequestWithContext(tb.Context(), http.MethodGet, "http://"+addr+path, nil)
	require.NoError(tb, err)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	require.NoError(tb, err)

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	require.NoError(tb, err)

	return reply{status: resp.StatusCode, body: string(body), cookies: resp.Cookies()}
}
