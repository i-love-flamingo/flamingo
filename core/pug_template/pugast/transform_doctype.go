package pugast

import "fmt"

// Render renders the doctype
func (d *Doctype) Render(p *PugAst, depth int) (string, bool) {
	p.Doctype = d.Val
	return fmt.Sprintf("<!DOCTYPE %s>\n", d.Val), false
}
