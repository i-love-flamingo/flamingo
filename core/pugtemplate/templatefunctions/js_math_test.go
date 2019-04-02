package templatefunctions

import (
	"testing"

	"go.aoe.com/flamingo/framework/template"

	"github.com/stretchr/testify/assert"
)

func TestJsMath(t *testing.T) {
	var jsMath template.Function = new(JsMath)

	assert.Equal(t, "Math", jsMath.Name())

	math := jsMath.Func().(func() Math)()

	// equal
	assert.Equal(t, 1., math.Min(1, int64(2), 3.))
	assert.Equal(t, 3., math.Max(1, int64(2), 3.))

	// ceil
	assert.Equal(t, 2, math.Ceil(2))
	assert.Equal(t, 2, math.Ceil(int64(2)))
	assert.Equal(t, 2, math.Ceil(2.))
	assert.Equal(t, 3, math.Ceil(2.4))
	assert.Equal(t, 3, math.Ceil(2.5))

	// trunc
	assert.Equal(t, 2, math.Trunc(2))
	assert.Equal(t, 2, math.Trunc(int64(2)))
	assert.Equal(t, 2, math.Trunc(2.))
	assert.Equal(t, 2, math.Trunc(2.1))

	// round
	assert.Equal(t, 2, math.Round(2))
	assert.Equal(t, 2, math.Round(int64(2)))
	assert.Equal(t, 2, math.Round(2.))
	assert.Equal(t, 2, math.Round(2.1))
	assert.Equal(t, 2, math.Round(2.4))
	assert.Equal(t, 3, math.Round(2.5))
	assert.Equal(t, 3, math.Round(2.9))
}
