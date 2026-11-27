package responder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload_Headers(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootResponder(t)

	resp := app.Get(t, "/download")

	assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
	assert.Equal(t, "attachment; filename=list.csv", resp.Header.Get("Content-Disposition"))
	assert.Equal(t, "a,b\n", resp.Body)
}
