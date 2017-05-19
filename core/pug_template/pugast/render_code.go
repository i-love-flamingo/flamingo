package pugast

// Render renders a code block
func (c *Code) Render(p *PugAst, depth int) (string, bool) {
	p.rawmode = !c.MustEscape
	return p.JsExpr(c.Val, true, true), *c.IsInline
}
