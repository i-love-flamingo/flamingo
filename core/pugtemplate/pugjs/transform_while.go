package pugjs

import (
	"bytes"
	"fmt"
)

// Render renders the loop, with obj or key+obj index
func (w *While) Render(p *renderState, wr *bytes.Buffer) error {
	fmt.Fprintf(wr, "{{ range %s -}}", p.JsExpr(w.Test, false, false))

	if err := w.Block.Render(p, wr); err != nil {
		return err
	}

	wr.WriteString("{{ end -}}")

	return nil
}
