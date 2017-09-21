package pugjs

import (
	"bytes"
	"fmt"
)

const casedefault = "default"

// Render a case node
func (c *Case) Render(s *renderState, wr *bytes.Buffer) error {
	var doElse string
	var elseBranch *When

	for _, node := range c.Block.Nodes {
		if node.(*When).Expr == casedefault {
			elseBranch = node.(*When)
		} else {
			fmt.Fprintf(wr, `{{- %sif __op__eql %s %s }}`, doElse, s.JsExpr(string(c.Expr), false, false), s.JsExpr(string(node.(*When).Expr), false, false))
			doElse = "else "
			if err := node.Render(s, wr); err != nil {
				return err
			}
		}
	}

	if elseBranch != nil {
		wr.WriteString(`{{- else }}`)
		if err := elseBranch.Render(s, wr); err != nil {
			return err
		}
	}
	wr.WriteString(`{{- end }}`)

	return nil
}

// Render a when node
func (w *When) Render(s *renderState, wr *bytes.Buffer) error {
	return w.Block.Render(s, wr)
}
