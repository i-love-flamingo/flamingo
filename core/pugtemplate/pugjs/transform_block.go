package pugjs

import "bytes"

// Render renders a Block, and intends every sub-block if necessary
func (b *Block) Render(s *renderState, wr *bytes.Buffer, depth int) error {
	for _, n := range b.Nodes {
		err := n.Render(s, wr, depth)
		if err != nil {
			return err
		}
	}
	return nil
}
