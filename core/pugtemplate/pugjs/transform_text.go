package pugjs

import "bytes"

// Render a text node
func (t *Text) Render(p *renderState, wr *bytes.Buffer) error {
	wr.WriteString(t.Val)
	return nil
}
