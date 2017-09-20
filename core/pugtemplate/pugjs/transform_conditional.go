package pugjs

import "bytes"

// Render renders a conditional via `if`
func (c *Conditional) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	wr.WriteString(`{{ if ` + p.JsExpr(string(c.Test), false, false) + ` -}}`)
	if err := c.Consequent.Render(p, wr, depth); err != nil {
		return err
	}

	if c.Alternate != nil {
		wr.WriteString(`{{ else -}}`)
		if err := c.Alternate.Render(p, wr, depth); err != nil {
			return err
		}
	}
	wr.WriteString(`{{ end -}}`)
	return nil
}
