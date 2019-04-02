package pugjs

import "bytes"

// Render renders a Block, and intends every sub-block if necessary
func (b *Block) Render(s *renderState, wr *bytes.Buffer) error {
	for _, n := range b.Nodes {
		err := n.Render(s, wr)
		if err != nil {
			return err
		}
	}
	return nil
}
