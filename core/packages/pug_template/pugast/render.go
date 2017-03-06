package pugast

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto/ast"
)

// TokenToTemplate gets named Template from Token
func (p *PugAst) TokenToTemplate(name string, t *Token) *template.Template {
	tpl := template.
		New(name).
		Funcs(FuncMap).
		Funcs(p.FuncMap).
		Option("missingkey=error")

	tc := p.render(t, "", nil)
	tpl, err := tpl.Parse(tc)
	p.TplCode[name] = tc

	if err != nil {
		e := err.Error() + "\n"
		for i, l := range strings.Split(tc, "\n") {
			e += fmt.Sprintf("%03d: %s\n", i+1, l)
		}
		panic(e)
	}

	return tpl
}

// args returns formatted keys and value string for given attributes from AST
func (p *PugAst) args(attrs []*Attr, andattributes bool) string {
	if len(attrs) == 0 {
		return ""
	}

	a := make(map[string]string)
	for _, attr := range attrs {
		p.rawmode = !attr.MustEscape
		a[attr.Name] += ` ` + p.JsExpr(attr.Val, true, false)
	}

	result := ""

	for k, v := range a {
		var aa string
		if andattributes {
			aa = ` {{index $__andattributes "` + k + `"}}`
		}
		if len(strings.TrimSpace(v)) > 0 {
			result += ` ` + k + `="` + strings.TrimSpace(v) + aa + `"`
		}
	}
	return result
}

// ifmt returns string for token, prefixed with newline if not inline token
func ifmt(t *Token, pre, buf string) string {
	if (t.IsInline != nil) && !*t.IsInline {
		return "\n" + pre + buf
	}
	return buf
}

const depth = "  "

// render returns formatted html
func (p *PugAst) render(parent *Token, pre string, mixinblock *Token) string {
	var buf string

	for _, t := range parent.Nodes {
		p.rawmode = !t.MustEscape

		switch t.Type {
		case "Extends", "RawInclude":
			if t.File.Path[0] == '/' {
				buf += p.render(p.Parse("jade"+t.File.Path+".jade"), pre, mixinblock)
			} else {
				buf += p.render(p.Parse(filepath.Dir(parent.Filename)+"/"+t.File.Path+".jade"), pre, mixinblock)
			}

		case "NamedBlock":
			switch t.Mode {
			case "replace":
				//if _, ok := blocks[t.Name]; !ok {
				//	buf += "\n" + pre + fmt.Sprintf("{{ pug_template \"%s\" . }}", t.Name)
				//}
				//blocks[t.Name] = p.render(t, pre+depth, mixinblock)
				buf += p.render(t, pre, mixinblock)
			//case "append":
			//	blocks[t.Name] += p.render(t, pre+depth, mixinblock)
			//case "prepend":
			//	blocks[t.Name] = p.render(t, pre+depth, mixinblock) + blocks[t.Name]
			default:
				panic(t.Mode)
			}

		case "Doctype":
			buf += fmt.Sprintf("<!DOCTYPE %s>", t.Val)

		case "Tag":
			if t.SelfClosing {
				buf += ifmt(t, pre, fmt.Sprintf("<%s />", t.Name))
			} else {
				andattrs := ""
				if len(t.AttributeBlocks) > 0 {
					buf += ifmt(t, pre, `{{$__andattributes := $`+t.AttributeBlocks[0]+`}}`)
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
				if (t.IsInline == nil || *t.IsInline) || (t.Block.IsInline == nil || *t.Block.IsInline) || len(t.Block.Nodes) == 0 {
					buf += ifmt(t, pre, fmt.Sprintf("<%s%s%s>", t.Name, p.args(t.Attrs, len(t.AttributeBlocks) > 0), andattrs))
					if t.Name != "meta" && t.Name != "link" {
						buf += fmt.Sprintf("%s</%s>", p.render(t.Block, "", mixinblock), t.Name)
					}
				} else if len(t.Block.Nodes) == 1 && (t.Block.Nodes[0].IsInline == nil || *t.Block.Nodes[0].IsInline) {
					buf += ifmt(t, pre, fmt.Sprintf("<%s%s%s>%s</%s>", t.Name, p.args(t.Attrs, len(t.AttributeBlocks) > 0), andattrs, p.render(t.Block, "", mixinblock), t.Name))
				} else {
					buf += ifmt(t, pre+depth, fmt.Sprintf("<%s%s%s>", t.Name, p.args(t.Attrs, len(t.AttributeBlocks) > 0), andattrs))
					buf += p.render(t.Block, pre+depth, mixinblock)
					buf += ifmt(t, pre+depth, fmt.Sprintf("</%s>", t.Name))
				}
			}

		/*
			case "InterpolatedTag":
				name := p.JsExpr(t.Expr, false, false)
				if t.SelfClosing {
					buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s ""}}/>`, name))
				} else {
					if t.IsInline || len(t.Block.Nodes) == 0 {
						buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s ""}}%s>%s</%s>`, name, p.args(t.Attrs, len(t.AttributeBlocks) > 0), p.render(t.Block, "", mixinblock), name))
					} else {
						buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s ""}}%s>`, name, p.args(t.Attrs, len(t.AttributeBlocks) > 0)))
						buf += p.render(t.Block, pre+depth, mixinblock)
						buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s "/"}}>`, name))
					}
				}
		*/

		case "Code":
			buf += ifmt(t, pre, p.JsExpr(t.Val, true, true))

		case "Mixin":
			if t.Call {
				if p.mixin[t.Name] == nil {
					panic("UNKNOWN MIXIN " + t.Name)
				}

				buf += "\n" + pre + "<!-- MIXIN " + t.Name + " -->"

				buf += "\n" + pre + `{{ $attributes := __op__map `
				for _, a := range t.Attrs {
					buf += ` "` + a.Name + `" ` + p.JsExpr(a.Val, false, false)
				}
				buf += " }}"

				callargs := strings.Split(p.mixin[t.Name].Args, ",")

				if len(callargs) == 1 {
					arg := p.JsExpr(t.Args, false, false)
					if len(arg) == 0 {
						arg = `null`
					}
					buf += "\n" + pre + fmt.Sprintf("{{ $%s := %s }}", callargs[0], arg)
					p.knownVar[callargs[0]] = true
				} else if len(callargs) > 1 {
					buf += "\n" + pre + fmt.Sprintf("{{ $__args__ := %s }}", p.JsExpr(`[`+t.Args+`]`, false, false))

					lenCaptured := 1
					if seq, ok := FuncToStatements(t.Args)[0].(*ast.ReturnStatement).Argument.(*ast.SequenceExpression); ok {
						lenCaptured = len(seq.Sequence)
					}

					for ci, ca := range callargs {
						ca = strings.TrimSpace(ca)
						p.knownVar[ca] = true
						if ci < lenCaptured {
							buf += "\n" + pre + fmt.Sprintf("{{ $%s := index $__args__ %d }}", ca, ci)
						} else {
							buf += "\n" + pre + fmt.Sprintf("{{ $%s := null }}", ca)
						}
					}
				}

				//buf += p.render(p.mixin[t.Name].Block, pre, t.Block)
				buf += "\n" + pre + "<!-- END MIXIN " + t.Name + " -->"
			} else {
				p.mixin[t.Name] = t
			}

		case "MixinBlock":
			if mixinblock != nil {
				buf += p.render(mixinblock, pre, mixinblock)
			}

		case "Comment", "BlockComment":
			buf += pre + "<!-- " + t.Val + " -->\n"

		case "Each":
			p.knownVar[t.Val] = true
			if t.Key != "" {
				p.knownVar[t.Key] = true
				buf += "\n" + pre + fmt.Sprintf("{{range $%s, $%s := %s}}", t.Key, t.Val, p.JsExpr(t.Obj, false, false))
			} else {
				buf += "\n" + pre + fmt.Sprintf("{{range $%s := %s}}", t.Val, p.JsExpr(t.Obj, false, false))
			}
			buf += p.render(t.Block, pre, mixinblock)
			buf += "\n" + pre + "{{end}}"

		case "Text":
			buf += t.Val

		case "Conditional":
			buf += "\n" + pre + fmt.Sprintf("{{if %s}}", p.JsExpr(t.Test, false, false))
			buf += p.render(t.Consequent, pre, mixinblock)
			if t.Alternate != nil {
				buf += "\n" + pre + "{{else}}"
				buf += p.render(t.Alternate, pre, mixinblock)
			}
			buf += "\n" + pre + "{{end}}"

		case "Block":
			buf += p.render(t, pre, mixinblock)

		default:
			panic(t.Type)
		}
	}

	return buf
}
