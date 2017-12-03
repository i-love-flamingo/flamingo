package pugjs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach_Render(t *testing.T) {
	var buffer = new(bytes.Buffer)
	var node = new(Each)
	var renderState = new(renderState)

	node.Val = "foo"
	node.Obj = "bar"

	assert.NoError(t, node.Render(renderState, buffer))
	assert.Equal(t, "{{ range $foo := $bar -}}{{ end -}}", buffer.String())

	buffer.Reset()

	node.Key = "key"

	assert.NoError(t, node.Render(renderState, buffer))
	assert.Equal(t, "{{ range $key, $foo := $bar -}}{{ end -}}", buffer.String())

	buffer.Reset()
}
