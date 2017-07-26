package pugjs

import (
	"bytes"
	"fmt"
)

// Render renders the loop, with obj or key+obj index
func (e *Each) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	if e.Key != "" {
		fmt.Fprintf(wr, "{{ range $%s, $%s := %s -}}", e.Key, e.Val, p.JsExpr(string(e.Obj), false, false))
	} else {
		fmt.Fprintf(wr, "{{ range $%s := %s -}}", e.Val, p.JsExpr(string(e.Obj), false, false))
	}
	if err := e.Block.Render(p, wr, depth); err != nil {
		return err
	}
	wr.WriteString("{{ end -}}")

	return nil
}
