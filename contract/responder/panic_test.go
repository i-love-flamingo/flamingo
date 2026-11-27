package responder_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestPanic_Gives500AndLogsStack(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootResponder(t)

	resp := app.Get(t, "/panic")

	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.Contains(t, resp.Body, "error-with-code")

	var found bool

	for _, entry := range app.Logs() {
		if entry.Level == zapcore.ErrorLevel &&
			strings.Contains(entry.Message, "catalog-secret-panic") &&
			strings.Contains(entry.Message, "testapp/modules/catalog.") {
			found = true
		}
	}

	assert.True(t, found, "no ERROR log entry with the panic text and the stack of the action")
}
