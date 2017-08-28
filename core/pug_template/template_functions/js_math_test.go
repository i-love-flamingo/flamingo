package template_functions

import (
	"flamingo/framework/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsMath(t *testing.T) {
	var jsMath template.Function = new(JsMath)

	assert.Equal(t, "Math", jsMath.Name())

	math := jsMath.Func().(func() Math)()

	assert.Equal(t, 1., math.Min(1, int64(2), 3.))
	assert.Equal(t, 3., math.Max(1, int64(2), 3.))

	assert.Equal(t, 2, math.Ceil(2))
	assert.Equal(t, 2, math.Ceil(int64(2)))
	assert.Equal(t, 2, math.Ceil(2.))
	assert.Equal(t, 3, math.Ceil(2.4))
	assert.Equal(t, 3, math.Ceil(2.5))
}
