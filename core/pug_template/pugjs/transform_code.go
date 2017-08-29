package pugjs

import "bytes"

// Render renders a code block
func (c *Code) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	p.rawmode = !c.MustEscape
	wr.WriteString(p.JsExpr(c.Val, true, true))
	return nil
}
