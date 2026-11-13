package config_test

import (
	"bytes"
	"log"
	"log/slog"
	"testing"
)

// captureSlog sends the default slog logger to a buffer for the rest of the
// test and returns that buffer. slog.SetDefault also redirects the log
// package, so cleanup restores the slog logger and the log writer and flags.
// Tests that use it must not run in parallel.
func captureSlog(t *testing.T) *bytes.Buffer {
	t.Helper()

	var output bytes.Buffer

	previousLogger, previousWriter, previousFlags := slog.Default(), log.Writer(), log.Flags()

	slog.SetDefault(slog.New(slog.NewTextHandler(&output, nil))) //nolint:forbidigo // The one place that swaps the default logger and restores it.
	t.Cleanup(func() {
		slog.SetDefault(previousLogger) //nolint:forbidigo // Restores the logger saved above.
		log.SetOutput(previousWriter)
		log.SetFlags(previousFlags)
	})

	return &output
}
