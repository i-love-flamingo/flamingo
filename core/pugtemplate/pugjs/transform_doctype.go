package pugjs

import (
	"bytes"
	"fmt"
)

// Render renders the doctype
func (d *Doctype) Render(p *renderState, wr *bytes.Buffer) error {
	p.doctype = d.Val
	fmt.Fprintf(wr, "<!DOCTYPE %s>\n", d.Val)
	return nil
}
