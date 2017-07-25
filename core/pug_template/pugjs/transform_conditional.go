package pugjs

// Render renders a conditional via `if`
func (c *Conditional) Render(p *renderState, depth int) (string, bool) {
	buf := `{{if ` + p.JsExpr(string(c.Test), false, false) + `}}`
	b, _ := c.Consequent.Render(p, depth)
	buf += b

	if c.Alternate != nil {
		buf += `{{else}}`
		b, _ := c.Alternate.Render(p, depth)
		buf += b
	}
	buf += `{{end}}`
	return buf, false
}
