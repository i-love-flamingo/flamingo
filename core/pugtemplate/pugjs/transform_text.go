package pugjs

import (
	"bytes"
	"strings"
)

// Render a text node
func (t *Text) Render(p *renderState, wr *bytes.Buffer) error {
	t.Val = strings.Replace(t.Val, "{{", `--{{--`, -1)
	t.Val = strings.Replace(t.Val, "}}", `--}}--`, -1)
	t.Val = strings.Replace(t.Val, "--{{--", `{{"{{"}}`, -1)
	t.Val = strings.Replace(t.Val, "--}}--", `{{"}}"}}`, -1)
	wr.WriteString(t.Val)
	return nil
}
