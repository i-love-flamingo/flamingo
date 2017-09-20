package pugjs

import (
	"bytes"
	"fmt"
	"strings"
)

// Render a tag
func (t *Tag) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	var _subblock = new(bytes.Buffer)
	if err := t.Block.Render(p, _subblock, depth+1); err != nil {
		return err
	}

	var attrs string
	if len(t.AttributeBlocks) > 0 || len(t.Attrs) > 0 {
		attrs = `{{ __attrs `
		for _, attr := range t.Attrs {
			if attr.MustEscape {
				attrs += fmt.Sprintf(`(__attr %q %s %t) `, attr.Name, p.JsExpr(string(attr.Val), false, false), attr.MustEscape)
			} else {
				attrs += fmt.Sprintf(`(__attr %q %q %t) `, attr.Name, p.JsExpr(string(attr.Val), false, false), attr.MustEscape)
			}
		}
		for _, ab := range t.AttributeBlocks {
			attrs += fmt.Sprintf(`(__and_attrs $%s)`, ab)
		}
		attrs += ` }}`
	}

	switch {
	case t.SelfClosing:
		fmt.Fprintf(wr, `<%s%s>`, t.Name, attrs)

	case t.Name == "script" && strings.Index(_subblock.String(), "\n") > -1:
		fmt.Fprintf(wr, "<%s%s>\n%s\n</%s>", t.Name, attrs, _subblock.String(), t.Name)

	case !t.Block.Inline() && p.debug:
		fmt.Fprintf(wr, "<%s%s>     {{- \"\" -}}\n%s     {{- \"\" -}}\n</%s>", t.Name, attrs, _subblock.String(), t.Name)

	default:
		fmt.Fprintf(wr, `<%s%s>%s</%s>`, t.Name, attrs, _subblock.String(), t.Name)
	}

	if !t.Inline() && p.debug {
		wr.WriteString("     {{- \"\" -}}\n")
	}

	return nil
}
