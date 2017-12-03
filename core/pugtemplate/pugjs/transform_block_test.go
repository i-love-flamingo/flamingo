package pugjs

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Render(t *testing.T) {
	var buffer = new(bytes.Buffer)
	var block = new(Block)
	var renderState = new(renderState)

	mock1, mock2, mock3 := new(MockNode), new(MockNode), new(MockNode)

	block.Nodes = []Node{
		mock1,
		mock2,
		mock3,
		mock1,
		mock1,
	}

	mock1.On("Render", renderState, buffer).Times(3).Return(nil)
	mock2.On("Render", renderState, buffer).Once().Return(nil)
	mock3.On("Render", renderState, buffer).Once().Return(nil)

	assert.NoError(t, block.Render(renderState, buffer))

	mock1.AssertExpectations(t)
	mock2.AssertExpectations(t)
	mock3.AssertExpectations(t)

	mock1.On("Render", renderState, buffer).Times(2).Return(nil)
	mock1.On("Render", renderState, buffer).Times(1).Return(errors.New("err"))
	block.Nodes = []Node{mock1, mock1, mock1}
	assert.Error(t, block.Render(renderState, buffer))

	mock1.AssertExpectations(t)
}
