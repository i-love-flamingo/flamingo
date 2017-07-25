package pugjs

import "fmt"

// Render renders the loop, with obj or key+obj index
func (e *Each) Render(p *renderState, depth int) (string, bool) {
	buf := ""

	if e.Key != "" {
		buf += fmt.Sprintf("{{range $%s, $%s := %s}}", e.Key, e.Val, p.JsExpr(string(e.Obj), false, false))
	} else {
		buf += fmt.Sprintf("{{range $%s := %s}}", e.Val, p.JsExpr(string(e.Obj), false, false))
	}
	b, _ := e.Block.Render(p, depth)
	buf += b
	buf += "{{end}}"

	return buf, false
}
