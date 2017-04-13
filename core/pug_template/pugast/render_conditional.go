package pugast

// Render renders a conditional via `if`
func (c *Conditional) Render(p *PugAst, depth int) (string, bool) {
	buf := `{{if ` + p.JsExpr(string(c.Test), false, false) + `}}`
	b, _ := c.Consequent.Render(p, depth)
	buf += b

	if len(c.Alternate.Nodes) > 0 {
		buf += `{{else}}`
		b, _ := c.Alternate.Render(p, depth)
		buf += b
	}
	buf += `{{end}}`
	return buf, false
}
