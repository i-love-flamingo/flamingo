package pugast

import "strings"

// Render renders a code block
func (c *Code) Render(p *PugAst, depth int) (string, bool) {
	p.rawmode = !c.MustEscape
	return strings.Replace(p.JsExpr(c.Val, true, true), "}}{{", "}}\n{{", -1), *c.IsInline
}
