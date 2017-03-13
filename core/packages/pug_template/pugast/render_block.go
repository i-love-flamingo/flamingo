package pugast

import "strings"

// Render renders a Block, and intends every sub-block if necessary
func (b *Block) Render(p *PugAst, depth int) (res string, isinline bool) {
	prefix := strings.Repeat("    ", depth)
	isinline = true
	for _, n := range b.Nodes {
		r, wasinline := n.Render(p, depth)
		if !wasinline {
			isinline = false
			res += "\n" + prefix
		}
		res += r
	}
	// remove double-indentation
	res = strings.Replace(res, "\n"+prefix+"\n"+prefix, "\n"+prefix, -1)
	return
}
