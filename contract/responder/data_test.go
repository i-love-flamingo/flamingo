package responder_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData_StatusContentTypeAndJSON(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootResponder(t)

	resp := app.Get(t, "/product/1")

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "{\"id\":\"1\"}\n", resp.Body)
}
