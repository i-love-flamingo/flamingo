package pugast

import (
	"fmt"
	"strings"
)

// args returns formatted keys and value string for given attributes from AST
func (t *Tag) args(p *PugAst, attrs []Attribute, andattributes bool) string {
	if len(attrs) == 0 {
		return ""
	}

	a := make(map[string]string)
	for _, attr := range attrs {
		p.rawmode = !attr.MustEscape
		a[attr.Name] += ` ` + p.JsExpr(string(attr.Val), true, false)
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
		k, v := attr.Name, a[attr.Name]
		if andattributes {
			aa = ` {{index $__andattributes "` + k + `"}}`
		}
		if len(strings.TrimSpace(v)) > 0 {
			result += ` ` + k + `="` + strings.TrimSpace(v) + aa + `"`
		}
	}

	return result
}

// Render a tag
func (t *Tag) Render(p *PugAst, depth int) (res string, isinline bool) {
	isinline = *t.IsInline
	prefix := strings.Repeat("    ", depth)

	subblock, wasinline := t.Block.Render(p, depth+1)

	andattrs := ""
	if len(t.AttributeBlocks) > 0 {
		res += `{{$__andattributes := $` + string(t.AttributeBlocks[0]) + `}}`
		knownaa := make(map[string]bool)
		for _, e := range t.Attrs {
			if len(e.Val) > 0 {
				knownaa[e.Name] = true
			}
		}
		for e := range knownaa {
			andattrs += ` "` + e + `"`
		}
		andattrs = ` {{__add_andattributes $__andattributes` + andattrs + `}}`
	}

	switch {
	case t.Name == "link", t.Name == "meta":
		res += fmt.Sprintf(`<%s%s%s>`, t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs)

	case t.SelfClosing:
		res += fmt.Sprintf(`<%s%s%s/>`, t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs)

	case t.IsInline != nil && !*t.IsInline:
		if !wasinline {
			subblock = subblock + "\n" + prefix
		}
		res += fmt.Sprintf("<%s%s%s>%s</%s>", t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs, subblock, t.Name)

	default:
		res += fmt.Sprintf(`<%s%s%s>%s</%s>`, t.Name, t.args(p, t.Attrs, len(t.AttributeBlocks) > 0), andattrs, subblock, t.Name)
	}

	return
}
