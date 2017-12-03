package pugjs

import (
	"bytes"
)

// Render a text node
func (t *Text) Render(p *renderState, wr *bytes.Buffer) error {
	_, err := wr.WriteString(t.Val)
	return err
}
