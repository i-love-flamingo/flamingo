package pugjs

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type (
	// OnRenderHTMLBlockEvent is an event which is called when a new block is going to be rendered
	OnRenderHTMLBlockEvent struct {
		BlockName string
		Buffer    *bytes.Buffer
	}
)

// Render an interpolated tag
func (it *InterpolatedTag) Render(p *renderState, wr *bytes.Buffer) error {
	return it.CommonTag.render(p.JsExpr(it.Expr, true, false), p, wr)
}

// Render a tag
func (t *Tag) Render(p *renderState, wr *bytes.Buffer) error {
	return t.CommonTag.render(t.Name, p, wr)
}

func (ct *CommonTag) render(name string, p *renderState, wr *bytes.Buffer) error {
	var subblock = new(bytes.Buffer)
	if err := ct.Block.Render(p, subblock); err != nil {
		return err
	}

	additional := new(bytes.Buffer)
	p.eventRouter.Dispatch(context.Background(), &OnRenderHTMLBlockEvent{name, additional})

	var attrs string
	if len(ct.AttributeBlocks) > 0 || len(ct.Attrs) > 0 {
		attrs = `{{ __attrs `
		for _, attr := range ct.Attrs {
			if attr.MustEscape {
				attrs += fmt.Sprintf(`(__attr %q %s %t) `, attr.Name, p.JsExpr(attr.Val, false, false), attr.MustEscape)
			} else {
				attrs += fmt.Sprintf(`(__attr %q %q %t) `, attr.Name, p.JsExpr(attr.Val, false, false), attr.MustEscape)
			}
		}
		for _, ab := range ct.AttributeBlocks {
			attrs += fmt.Sprintf(`(__and_attrs $%s)`, ab)
		}
		attrs += ` }}`
	}

	switch {
	case ct.SelfClosing:
		fmt.Fprintf(wr, `<%s%s>`, name, attrs)

	case name == "script" && strings.Index(subblock.String(), "\n") > -1:
		fmt.Fprintf(wr, "<%s%s>\n%s\n</%s>", name, attrs, additional.String()+subblock.String(), name)

	case !ct.Block.Inline() && p.debug:
		fmt.Fprintf(wr, "<%s%s>     {{- \"\" -}}\n%s     {{- \"\" -}}\n</%s>", name, attrs, additional.String()+subblock.String(), name)

	default:
		fmt.Fprintf(wr, `<%s%s>%s</%s>`, name, attrs, additional.String()+subblock.String(), name)
	}

	if !ct.Inline() && p.debug {
		wr.WriteString("     {{- \"\" -}}\n")
	}

	return nil
}
