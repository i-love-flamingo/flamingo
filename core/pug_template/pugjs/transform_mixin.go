package pugjs

import (
	"bytes"
	"fmt"
	"strings"
)

// Render renders the mixin, either it's call or it's definition
func (m *Mixin) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	if m.Call {
		return m.renderCall(p, wr, depth)
	}

	return m.renderDefinition(p, wr, depth)
}

func (m *Mixin) renderDefinition(p *renderState, wr *bytes.Buffer, depth int) error {
	if p.mixin[string(m.Name)] != "" {
		return nil
	}

	callargs := strings.Split(m.Args, ",")
	attrpart := ""

	for ci, ca := range callargs {
		ca = strings.TrimSpace(ca)
		attrpart += fmt.Sprintf(
			"{{- $%s := __tryindex $__args__ %d -}}",
			ca,
			ci,
		)
	}

	var subblock = new(bytes.Buffer)

	if err := m.Block.Render(p, subblock, depth); err != nil {
		return err
	}

	prefix := strings.Repeat("  ", depth)
	p.mixin[string(m.Name)] = fmt.Sprintf(`
{{- define "mixin_%s" }}
{{- $attributes := (__tryindex . 1) }}
{{- $__args__ := (__tryindex . 0) }}
{{- $block := (__tryindex . 2) }}
%s%s
%s%s
{{- end }}`, m.Name, prefix, attrpart, prefix, subblock.String())
	p.mixinorder = append(p.mixinorder, string(m.Name))
	return nil
}

func (m *Mixin) renderCall(p *renderState, wr *bytes.Buffer, depth int) error {
	attributes := `__op__map_params `
	for _, a := range m.Attrs {
		attributes += ` "` + a.Name + `" ` + p.JsExpr(string(a.Val), false, false)
	}
	var subblock = new(bytes.Buffer)
	if err := m.Block.Render(p, subblock, depth); err != nil {
		return err
	}
	if len(subblock.String()) > 0 {
		prefix := strings.Repeat("  ", depth)
		blockname := fmt.Sprintf("block_%s_%d", m.Name, p.mixincounter)
		p.mixincounter++
		mixinblock := fmt.Sprintf(`
%s{{- define "%s" -}}
%s%s
%s{{- end -}}`, prefix, blockname, prefix, subblock.String(), prefix)
		p.mixinblocks = append(p.mixinblocks, mixinblock)
		fmt.Fprintf(wr, `{{ template "mixin_%s" (__op__array (%s) (%s) ("%s") ) }}`, m.Name, p.JsExpr(`[`+m.Args+`]`, false, false), attributes, blockname)
	} else {
		fmt.Fprintf(wr, `{{ template "mixin_%s" (__op__array (%s) (%s) (null) ) }}`, m.Name, p.JsExpr(`[`+m.Args+`]`, false, false), attributes)
	}
	return nil
}

// Render MixinBlock call
func (m *MixinBlock) Render(p *renderState, wr *bytes.Buffer, depth int) error {
	wr.WriteString(`{{- template $block -}}`)
	return nil
}
