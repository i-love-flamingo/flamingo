package pugjs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText_Render(t *testing.T) {
	var buffer = new(bytes.Buffer)
	var text = new(Text)

	text.Val = "test 123"

	assert.NoError(t, text.Render(new(renderState), buffer))
	assert.Equal(t, "test 123", buffer.String())
}
