package pugast

import "fmt"

// Render renders the doctype
func (d *Doctype) Render(p *PugAst, depth int) (string, bool) {
	return fmt.Sprintf("<!DOCTYPE %s>\n", d.Val), false
}
