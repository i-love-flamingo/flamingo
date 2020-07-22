package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyResults(t *testing.T) {
	recorder := httptest.NewRecorder()
	assert.NoError(t, new(Response).Apply(context.Background(), recorder))
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder = httptest.NewRecorder()
	assert.Error(t, new(RouteRedirectResponse).Apply(context.Background(), recorder))

	recorder = httptest.NewRecorder()
	assert.Error(t, new(URLRedirectResponse).Apply(context.Background(), recorder))

	recorder = httptest.NewRecorder()
	assert.NoError(t, new(DataResponse).Apply(context.Background(), recorder))
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder = httptest.NewRecorder()
	assert.NoError(t, new(RenderResponse).Apply(context.Background(), recorder))
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder = httptest.NewRecorder()
	assert.NoError(t, new(ServerErrorResponse).Apply(context.Background(), recorder))
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
