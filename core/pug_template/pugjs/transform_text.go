package pugjs

// Render a text node
func (t *Text) Render(p *renderState, depth int) (string, bool) {
	return t.Val, true
}
