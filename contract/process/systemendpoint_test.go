package process_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemEndpoint_Routes(t *testing.T) {
	t.Parallel()

	_, addrs := serve(t, testappMain, nil)

	for _, tc := range []struct {
		path string
		body string // text that the body contains
	}{
		{"/status/ping", "OK"},
		{"/status/healthcheck", `{"name":"ops","alive":true,"details":"ops up"}`},
		{"/metrics", "# HELP"},
		{"/version", "App version:"},
	} {
		got := get(t, addrs.System, tc.path)

		assert.Equal(t, 200, got.status, tc.path)
		assert.Contains(t, got.body, tc.body, tc.path)
	}
}
