package pugjs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoctype_Render(t *testing.T) {
	var buffer = new(bytes.Buffer)
	var node = new(Doctype)
	var renderState = new(renderState)

	node.Val = "html"

	assert.NoError(t, node.Render(renderState, buffer))
	assert.Equal(t, "<!DOCTYPE html>\n", buffer.String())
	assert.Equal(t, "html", renderState.doctype)
}
