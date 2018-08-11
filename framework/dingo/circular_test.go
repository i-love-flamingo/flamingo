package dingo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	circA struct {
		A *circA `inject:""`
		B *circB `inject:""`
	}

	circB struct {
		A *circA `inject:""`
		B *circB `inject:""`
	}
)

func TestDingoCircula(t *testing.T) {
	traceCircular = make([]circularTraceEntry, 0)
	defer func() {
		traceCircular = nil
	}()

	injector := NewInjector()
	assert.Panics(t, func() {
		_, ok := injector.GetInstance(new(circA)).(*circA)
		if !ok {
			t.Fail()
		}
	})
}
