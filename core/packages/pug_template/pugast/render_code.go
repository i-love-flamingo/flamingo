package pugast

// Render renders a code block
func (c *Code) Render(p *PugAst, depth int) (string, bool) {
	return p.JsExpr(c.Val, true, true), *c.IsInline
}
