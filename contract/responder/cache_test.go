package responder_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheHeaders_DirectiveAndNoCache(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootResponder(t)

	resp := app.Get(t, "/cached")

	assert.ElementsMatch(t, []string{"public", "max-age=60"}, strings.Split(resp.Header.Get("Cache-Control"), ", "))

	resp = app.Get(t, "/nocache")

	assert.Equal(t, "no-store", resp.Header.Get("Cache-Control"))
}
