package pugjs

import (
	"bytes"
	"strings"
)

// Render renders a code block
func (c *Code) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	p.rawmode = !c.MustEscape
	wr.WriteString(strings.Replace(p.JsExpr(c.Val, true, true), " }}{{", " -}}\n{{", -1))
	return nil
}
