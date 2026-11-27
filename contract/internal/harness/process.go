package harness

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	processTimeout = 30 * time.Second
	pollInterval   = 10 * time.Millisecond
	logTail        = 50
	maxLineSize    = 1 << 20
	appEnv         = "FLAMINGO_CONTRACT_APP" // names the app that the test binary runs instead of the tests
	unknownAppExit = 2                       // exit code for an app name that Main does not know
)

var (
	httpLine   = regexp.MustCompile(`^Starting HTTP Server(?: \(Prefixrouter\))? at (\S+?)(?: \.+)?$`)
	systemLine = regexp.MustCompile(`^Starting HTTP Server \(systemendpoint\) at (\S+)$`)

	// mains holds the apps that Main got.
	mains map[string]func()
)

// Main runs the tests, or the app that Process asks for with appEnv. Call it from TestMain.
// The apps are part of the test binary, so go test does not use a cached result after their code changed.
func Main(m *testing.M, apps map[string]func()) {
	mains = apps

	if name := os.Getenv(appEnv); name != "" {
		run, ok := apps[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "%s: unknown app %q\n", appEnv, name)
			os.Exit(unknownAppExit)
		}

		run()
		os.Exit(0)
	}

	os.Exit(m.Run())
}

// Options describes the process that Process starts.
type Options struct {
	Main string            // name of the app in the map given to Main
	Dir  string            // working directory of the app, "" for the test directory
	Env  map[string]string // added to the environment of the test
	Args []string          // flamingo flags, then the command: "serve", "--addr", "127.0.0.1:0"
}

// Addrs are the listen addresses that the process logs at the start.
type Addrs struct {
	HTTP   string // from "Starting HTTP Server at <addr> ....." or "Starting HTTP Server (Prefixrouter) at <addr>"
	System string // from "Starting HTTP Server (systemendpoint) at <addr>"
}

// Proc is a running process of an app.
type Proc struct {
	cmd *exec.Cmd

	mu    sync.Mutex
	lines []string

	done chan struct{} // closed after the end of stderr and the exit of the process
	exit int
}

// Process starts the test binary again as the app that opts.Main names, and kills it when the test ends.
// TestMain must call Main with that app. The app has its own dingo scopes, so the test can be parallel.
// If the test fails, the last lines of stderr go to the test log.
func Process(tb testing.TB, opts Options) *Proc {
	tb.Helper()

	if mains == nil {
		tb.Fatalf("process: TestMain does not call harness.Main")
	}

	if _, ok := mains[opts.Main]; !ok {
		tb.Fatalf("process: no app %q in the map given to harness.Main", opts.Main)
	}

	binary, err := os.Executable()
	if err != nil {
		tb.Fatalf("process: find the test binary: %v", err)
	}

	readFiles(tb, opts.Dir)

	cmd := exec.CommandContext(tb.Context(), binary, opts.Args...)
	cmd.Dir = opts.Dir
	// The CONTEXT and CONTEXTFILE of the shell must not change the config. opts.Env can set them again.
	cmd.Env = append(os.Environ(), "CONTEXT=", "CONTEXTFILE=", appEnv+"="+opts.Main)
	for name, value := range opts.Env {
		cmd.Env = append(cmd.Env, name+"="+value)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		tb.Fatalf("process: stderr pipe: %v", err)
	}

	proc := &Proc{cmd: cmd, done: make(chan struct{})}

	if err := cmd.Start(); err != nil {
		tb.Fatalf("process: start %s: %v", opts.Main, err)
	}

	go proc.read(stderr)

	tb.Cleanup(func() {
		proc.stop(tb)
	})

	return proc
}

// readFiles reads each file below dir. go test records the files that a test opens, so a changed config file
// of the app also invalidates a cached test result.
func readFiles(tb testing.TB, dir string) {
	tb.Helper()

	if dir == "" {
		dir = "."
	}

	fsys := os.DirFS(dir)

	err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.Type().IsRegular() {
			if _, err := fs.ReadFile(fsys, path); err != nil {
				tb.Fatalf("process: read file %s: %v", path, err)
			}
		}

		return nil
	})
	if err != nil {
		tb.Fatalf("process: walk dir %s: %v", dir, err)
	}
}

// read collects the lines of stderr, then waits for the exit of the process.
func (p *Proc) read(stderr io.Reader) {
	defer close(p.done)

	scanner := bufio.NewScanner(stderr)
	scanner.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), maxLineSize)

	for scanner.Scan() {
		p.mu.Lock()
		p.lines = append(p.lines, scanner.Text())
		p.mu.Unlock()
	}

	// Wait closes the pipe, so it must come after the last read.
	_ = p.cmd.Wait()

	p.exit = p.cmd.ProcessState.ExitCode()
}

func (p *Proc) stop(tb testing.TB) {
	tb.Helper()

	select {
	case <-p.done:
	default:
		_ = p.cmd.Process.Kill()

		select {
		case <-p.done:
		case <-time.After(processTimeout):
			tb.Errorf("process: still runs after kill")
		}
	}

	if tb.Failed() {
		lines := p.Lines()
		if len(lines) > logTail {
			lines = lines[len(lines)-logTail:]
		}

		tb.Logf("process: last %d lines of stderr:\n%s", len(lines), strings.Join(lines, "\n"))
	}
}

// WaitReady waits for the log lines of the HTTP server and of the system endpoint.
func (p *Proc) WaitReady(tb testing.TB) Addrs {
	tb.Helper()

	deadline := time.Now().Add(processTimeout)

	for {
		var addrs Addrs

		for _, entry := range p.Logs() {
			message, _ := entry["message"].(string)

			if match := httpLine.FindStringSubmatch(message); match != nil {
				addrs.HTTP = match[1]
			}

			if match := systemLine.FindStringSubmatch(message); match != nil {
				addrs.System = match[1]
			}
		}

		if addrs.HTTP != "" && addrs.System != "" {
			return addrs
		}

		select {
		case <-p.done:
			tb.Fatalf("process: ended before both servers started, exit code %d", p.exit)
		case <-time.After(pollInterval):
		}

		if time.Now().After(deadline) {
			tb.Fatalf("process: servers not started after %s, got %+v", processTimeout, addrs)
		}
	}
}

// Logs returns the lines of stderr that are JSON.
func (p *Proc) Logs() []map[string]any {
	var entries []map[string]any

	for _, line := range p.Lines() {
		var entry map[string]any

		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			entries = append(entries, entry)
		}
	}

	return entries
}

// Lines returns all lines of stderr up to now.
func (p *Proc) Lines() []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return slices.Clone(p.lines)
}

// Signal sends a signal to the process.
func (p *Proc) Signal(tb testing.TB, sig os.Signal) {
	tb.Helper()

	if err := p.cmd.Process.Signal(sig); err != nil {
		tb.Fatalf("process: send %v: %v", sig, err)
	}
}

// Wait waits for the end of the process and returns its exit code. It stops the test after 30 seconds.
func (p *Proc) Wait(tb testing.TB) int {
	tb.Helper()

	select {
	case <-p.done:
		return p.exit
	case <-time.After(processTimeout):
		tb.Fatalf("process: still runs after %s", processTimeout)

		return 0
	}
}
