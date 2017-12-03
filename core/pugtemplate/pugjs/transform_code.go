package pugjs

import "bytes"

// Render renders a code block
func (c *Code) Render(p *renderState, wr *bytes.Buffer) error {
	p.rawmode = !c.MustEscape
	_, err := wr.WriteString(p.JsExpr(JavaScriptExpression(c.Val), true, true))
	return err
}
