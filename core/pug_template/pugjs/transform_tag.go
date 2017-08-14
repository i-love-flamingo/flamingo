package pugjs

import (
	"bytes"
	"fmt"
	"strings"
)

// args returns formatted keys and value string for given attributes from AST
func (t *Tag) args(p *renderState, attrs []Attribute, andattributes bool) string {
	if len(attrs) == 0 {
		return ""
	}

	a := make(map[string]string)
	for _, attr := range attrs {
		p.rawmode = !attr.MustEscape
		if attr.Name == "style" {
			a[attr.Name] += ` ` + strings.Replace(p.JsExpr(string(attr.Val), true, false), `{{(__str `, `{{(sc `, -1)
		} else {
			a[attr.Name] += ` ` + strings.Replace(p.JsExpr(string(attr.Val), true, false), `{{(`, `{{html (`, -1)
		}
	}

	result := ""

	visited := make(map[string]bool)
	for _, attr := range attrs {
		if !visited[attr.Name] {
			visited[attr.Name] = true
		} else {
			continue
		}

		var aa string
		k, v := attr.Name, strings.TrimSpace(a[attr.Name])
		if andattributes {
			aa = `{{__str " " (index $__andattributes "` + k + `")}}`
		}
		if p.doctype == "html" && v == "true" {
			result += ` ` + k
		} else if p.doctype == "html" && v == "false" {
			// empty
		} else if len(v) > 0 {
			result += ` ` + k + `="` + v + aa + `"`
		}
	}

	return result
}

// Render a tag
func (t *Tag) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	var _subblock = new(bytes.Buffer)
	if err := t.Block.Render(p, _subblock, depth+1); err != nil {
		return err
	}
	subblock := _subblock.String()

	if strings.Index(subblock, "\n") > -1 {
		lines := strings.Split(subblock, "\n")
		for i, line := range lines {
			lines[i] = "  " + line
		}
		subblock = strings.Join(lines, "\n")
	}

	andattrs := ""
	if len(t.AttributeBlocks) > 0 {
		wr.WriteString(`{{$__andattributes := $` + string(t.AttributeBlocks[0]) + `}}`)
		knownaa := make(map[string]bool)
		for _, e := range t.Attrs {
			if len(e.Val) > 0 {
				knownaa[e.Name] = true
			}
		}
		for e := range knownaa {
			andattrs += ` "` + e + `"`
		}
		andattrs = `{{__add_andattributes $__andattributes` + andattrs + `}}`
	}

	switch {
	case t.Name == "link", t.Name == "meta":
		fmt.Fprintf(wr, `<%s%s%s>`, t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs)

	case t.SelfClosing:
		fmt.Fprintf(wr, `<%s%s%s/>`, t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs)

	case !t.Block.Inline() || (t.Name == "script" && strings.Index(subblock, "\n") > -1):
		fmt.Fprintf(wr, "<%s%s%s>\n%s\n</%s>", t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs, subblock, t.Name)

	default:
		fmt.Fprintf(wr, `<%s%s%s>%s</%s>`, t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs, subblock, t.Name)
	}

	if !t.Inline() {
		wr.WriteString("\n")
	}

	return nil
}
