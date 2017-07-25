package pugast

import "fmt"

const defaul = "default"

// Render a case node
func (c *Case) Render(p *PugAst, depth int) (string, bool) {

	buf := `{{if false}}`

	var els *When

	for _, node := range c.Block.Nodes {
		if node.(*When).Expr == defaul {
			els = node.(*When)
		} else {
			buf += fmt.Sprintf(`{{else if __op__eql %s %s}}`, p.JsExpr(string(c.Expr), false, false), p.JsExpr(string(node.(*When).Expr), false, false))
			b, _ := node.Render(p, depth+1)
			buf += b
		}
	}

	if els != nil {
		buf += `{{else}}`
		b, _ := els.Render(p, depth+1)
		buf += b
	}

	buf += `{{end}}`

	return buf, false
}

// Render a when node
func (w *When) Render(p *PugAst, depth int) (string, bool) {
	return w.Block.Render(p, depth)
}
