//go:build !race

package perf_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"flamingo.me/flamingo/v3/framework/web"
)

const (
	dataResultCeiling = 21
	dataResultRuns    = 100
)

func TestPerf_AllocsDataResult(t *testing.T) { //nolint:paralleltest // allocation counts need a test that is not parallel
	responder := new(web.Responder)

	got := minAllocs(t, dataResultRuns, func() {
		if err := responder.Data(map[string]string{"a": "b"}).Apply(context.Background(), httptest.NewRecorder()); err != nil {
			t.Error(err)
		}
	})

	assertAllocs(t, got, dataResultCeiling)
}
