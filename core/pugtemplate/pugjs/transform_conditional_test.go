package pugjs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditional_Render(t *testing.T) {
	var buffer = new(bytes.Buffer)
	var node = new(Conditional)
	var renderState = new(renderState)

	assert.Error(t, node.Render(renderState, buffer))

	buffer.Reset()
	node.Test = "1 == 2"
	node.Consequent = new(MockNode)
	node.Consequent.(*MockNode).On("Render", renderState, buffer).Once().Return(nil)

	assert.NoError(t, node.Render(renderState, buffer))
	assert.Equal(t, "{{ if (__op__eql 1 2) -}}{{ end -}}", buffer.String())

	node.Consequent.(*MockNode).AssertExpectations(t)

	buffer.Reset()
	node.Consequent.(*MockNode).On("Render", renderState, buffer).Once().Return(nil)
	node.Alternate = new(MockNode)
	node.Alternate.(*MockNode).On("Render", renderState, buffer).Once().Return(nil)

	assert.NoError(t, node.Render(renderState, buffer))
	assert.Equal(t, "{{ if (__op__eql 1 2) -}}{{ else -}}{{ end -}}", buffer.String())

	node.Consequent.(*MockNode).AssertExpectations(t)
	node.Alternate.(*MockNode).AssertExpectations(t)
}
