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
	prefix := strings.Repeat("    ", depth)

	callargs := strings.Split(m.Args, ",")
	attrpart := ""

	for ci, ca := range callargs {
		ca = strings.TrimSpace(ca)
		p.knownVar[ca] = true
		attrpart += fmt.Sprintf(
			"{{- $%s := tryindex $__args__ %d -}}",
			ca,
			ci,
		)
	}

	subblock, _ := m.Block.Render(p, depth)

	return fmt.Sprintf(`{{- define "mixin_%s" -}}
{{- $attributes := (index . 1) -}}
{{- $__args__ := (index . 0) -}}
%s
%s%s
%s
{{- end -}}
`, m.Name, attrpart, prefix, subblock, prefix)
}

func (m *Mixin) renderCall(p *PugAst, depth int) string {
	attributes := `__op__map `
	for _, a := range m.Attrs {
		attributes += ` "` + a.Name + `" ` + p.JsExpr(string(a.Val), false, false)
	}
	return fmt.Sprintf("{{ template \"mixin_%s\" (__op__array (%s) (%s) ) }}", m.Name, p.JsExpr(`[`+m.Args+`]`, false, false), attributes)
}
