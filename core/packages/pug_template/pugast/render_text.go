package pugast

// Render a text node
func (t *Text) Render(p *PugAst, depth int) (string, bool) {
	return t.Val, true
}
