package pugast

import (
	"fmt"
	"strings"
)

// Render renders the mixin, either it's call or it's definition
func (m *Mixin) Render(p *PugAst, depth int) (string, bool) {
	if m.Call {
		return m.renderCall(p, depth), false
	}

	return m.renderDefinition(p, depth), false
}

func (m *Mixin) renderDefinition(p *PugAst, depth int) string {
	if p.mixin[string(m.Name)] != "" {
		return ""
	}

	prefix := strings.Repeat("    ", depth)

	callargs := strings.Split(m.Args, ",")
	attrpart := ""

	for ci, ca := range callargs {
		ca = strings.TrimSpace(ca)
		attrpart += fmt.Sprintf(
			"{{- $%s := tryindex $__args__ %d -}}",
			ca,
			ci,
		)
	}

	subblock, _ := m.Block.Render(p, depth)

	p.mixin[string(m.Name)] = fmt.Sprintf(`{{- define "mixin_%s" -}}
{{- $attributes := (tryindex . 1) -}}
{{- $__args__ := (tryindex . 0) -}}
{{- $block := (tryindex . 2) -}}
%s
%s%s
%s
{{- end -}}
`, m.Name, attrpart, prefix, subblock, prefix)
	return ""
}

func (m *Mixin) renderCall(p *PugAst, depth int) string {
	attributes := `__op__map `
	for _, a := range m.Attrs {
		attributes += ` "` + a.Name + `" ` + p.JsExpr(string(a.Val), false, false)
	}
	block, _ := m.Block.Render(p, depth)
	if len(block) > 0 {
		blockname := fmt.Sprintf("block_%s_%d", m.Name, p.mixincounter)
		p.mixincounter++
		block = fmt.Sprintf(`
		{{- define "%s" -}}
		%s
		{{- end -}}
		`, blockname, block)
		p.mixinblocks = append(p.mixinblocks, block)
		return fmt.Sprintf(`{{ template "mixin_%s" (__op__array (%s) (%s) ("%s") ) }}`, m.Name, p.JsExpr(`[`+m.Args+`]`, false, false), attributes, blockname)
	}
	return fmt.Sprintf(`{{ template "mixin_%s" (__op__array (%s) (%s) (null) ) }}`, m.Name, p.JsExpr(`[`+m.Args+`]`, false, false), attributes)
}

func (m *MixinBlock) Render(p *PugAst, depth int) (string, bool) {
	return `{{ template $block }}`, false
}
