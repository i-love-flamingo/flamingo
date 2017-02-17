package node

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

var mixin map[string]*Token
var TplCode map[string]string

func init() {
	mixin = make(map[string]*Token)
}

func (p *PugAst) TokenToTemplate(name string, t *Token) *template.Template {
	TplCode = make(map[string]string)

	tpl := template.New(name).Funcs(FuncMap).Option("missingkey=error")

	tc := p.render(t, "", nil)
	tpl, err := tpl.Parse(tc)
	TplCode[name] = tc

	if err != nil {
		e := err.Error() + "\n"
		for i, l := range strings.Split(tc, "\n") {
			e += fmt.Sprintf("%03d: %s\n", i+1, l)
		}
		panic(e)
	}

	/*
		for name, block := range blocks {
			tpl, err = tpl.New(name).Parse(block)
			if err != nil {
				panic(err)
			}
		}
	*/

	return tpl
}

func args(attrs []*Attr) string {
	if len(attrs) == 0 {
		return ""
	}

	a := make(map[string]string)
	for _, attr := range attrs {
		a[attr.Name] += ` ` + JsExpr(attr.Val, true, false)
	}
	res := ""
	for k, v := range a {
		res += ` ` + k + `="` + strings.TrimSpace(v) + `"`
	}
	return res
}

func ifmt(t *Token, pre, buf string) string {
	if t.IsInline {
		return buf
	} else {
		return "\n" + pre + buf
	}
}

const depth = "  "

func (p *PugAst) render(parent *Token, pre string, mixinblock *Token) string {
	var buf string

	for _, t := range parent.Nodes {
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
				/*
					if _, ok := blocks[t.Name]; !ok {
						buf += "\n" + pre + fmt.Sprintf("{{ template \"%s\" . }}", t.Name)
					}
					blocks[t.Name] = p.render(t, pre+depth, mixinblock)
				*/
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
				if t.IsInline || len(t.Block.Nodes) == 0 {
					buf += ifmt(t, pre, fmt.Sprintf("<%s%s>%s</%s>", t.Name, args(t.Attrs), p.render(t.Block, "", mixinblock), t.Name))
				} else {
					buf += ifmt(t, pre, fmt.Sprintf("<%s%s>", t.Name, args(t.Attrs)))
					buf += p.render(t.Block, pre+depth, mixinblock)
					buf += ifmt(t, pre, fmt.Sprintf("</%s>", t.Name))
				}
			}

		case "InterpolatedTag":
			name := JsExpr(t.Expr, false, false)
			if t.SelfClosing {
				buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s ""}}/>`, name))
			} else {
				if t.IsInline || len(t.Block.Nodes) == 0 {
					buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s ""}}%s>%s</%s>`, name, args(t.Attrs), p.render(t.Block, "", mixinblock), name))
				} else {
					buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s ""}}%s>`, name, args(t.Attrs)))
					buf += p.render(t.Block, pre+depth, mixinblock)
					buf += ifmt(t, pre, fmt.Sprintf(`{{tagopen %s "/"}}>`, name))
				}
			}

		case "Code":
			buf += ifmt(t, pre, JsExpr(t.Val, true, true))

		case "Mixin":
			if t.Call {
				if mixin[t.Name] == nil {
					panic("UNKNOWN MIXIN " + t.Name)
				}

				buf += "\n" + pre + "<!-- MIXIN " + t.Name + " -->"

				buf += "\n" + pre + `{{ $attributes := __op__map `
				for _, a := range t.Attrs {
					buf += ` "` + a.Name + `" ` + JsExpr(a.Val, false, false)
				}
				buf += " }}"

				callargs := strings.Split(mixin[t.Name].Args, ",")

				if len(callargs) == 1 {
					arg := JsExpr(t.Args, false, false)
					if len(arg) == 0 {
						arg = `null`
					}
					buf += "\n" + pre + fmt.Sprintf("{{ $%s := %s }}", callargs[0], arg)
					known[callargs[0]] = true
				} else if len(callargs) > 1 {
					buf += "\n" + pre + fmt.Sprintf("{{ $__args__ := %s }}", JsExpr(t.Args, false, false))

					for ci, ca := range callargs {
						ca = strings.TrimSpace(ca)
						known[ca] = true
						buf += "\n" + pre + fmt.Sprintf("{{ $%s := index $__args__ %d }}", ca, ci)
					}
				}

				buf += p.render(mixin[t.Name].Block, pre, t.Block)
				buf += "\n" + pre + "<!-- END MIXIN " + t.Name + " -->"
			} else {
				mixin[t.Name] = t
			}

		case "MixinBlock":
			if mixinblock != nil {
				buf += p.render(mixinblock, pre, mixinblock)
			}

		case "Comment", "BlockComment":
			//buf += pre + "<!-- " + t.Val + " -->\n"

		case "Each":
			known[t.Val] = true
			if t.Key != "" {
				known[t.Key] = true
				buf += "\n" + pre + fmt.Sprintf("{{range $%s, $%s := %s}}", t.Key, t.Val, JsExpr(t.Obj, false, false))
			} else {
				buf += "\n" + pre + fmt.Sprintf("{{range $%s := %s}}", t.Val, JsExpr(t.Obj, false, false))
			}
			buf += p.render(t.Block, pre, mixinblock)
			buf += "\n" + pre + "{{end}}"

		case "Text":
			buf += ifmt(t, pre, t.Val)

		case "Conditional":
			buf += "\n" + pre + fmt.Sprintf("{{if %s}}", JsExpr(t.Test, false, false))
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
