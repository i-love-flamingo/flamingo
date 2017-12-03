package pugjs

import (
	"bytes"

	"github.com/pkg/errors"
)

// Render renders a conditional via `if`
func (c *Conditional) Render(p *renderState, wr *bytes.Buffer) error {
	wr.WriteString(`{{ if ` + p.JsExpr(c.Test, false, false) + ` -}}`)

	if c.Consequent == nil {
		return errors.New("can not render conditional without consequent")
	}

	if err := c.Consequent.Render(p, wr); err != nil {
		return err
	}

	if c.Alternate != nil {
		wr.WriteString(`{{ else -}}`)
		if err := c.Alternate.Render(p, wr); err != nil {
			return err
		}
	}
	wr.WriteString(`{{ end -}}`)
	return nil
}
