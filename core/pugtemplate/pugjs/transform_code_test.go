package pugjs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode_Render(t *testing.T) {
	var buffer = new(bytes.Buffer)
	var node = new(Code)

	node.Val = "var foo = 1"

	assert.NoError(t, node.Render(new(renderState), buffer))
	assert.Equal(t, "{{ $foo := 1 -}}", buffer.String())
}
