package pugast

import (
	"fmt"
	"html/template"
	"log"
	"strings"
)

// TokenToTemplate gets named Template from Token
func (p *PugAst) TokenToTemplate(name string, t *Token) *template.Template {
	tpl := template.
		New(name).
		Funcs(FuncMap).
		Funcs(p.FuncMap).
		Option("missingkey=error")

	tr := p.build(t)
	tc := ""

	for _, b := range tr {
		bla, _ := b.Render(p, 0)
		tc += bla
	}

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

func (p *PugAst) build(parent *Token) (res []Node) {
	if parent == nil {
		return
	}
	for _, t := range parent.Nodes {
		switch t.Type {
		case "Tag":
			tag := new(Tag)
			tag.Name = t.Name
			tag.IsInline = t.IsInline
			tag.SelfClosing = t.SelfClosing
			for _, a := range t.AttributeBlocks {
				tag.AttributeBlocks = append(tag.AttributeBlocks, JavaScriptExpression(a))
			}
			tag.Block = Block{Nodes: p.build(t.Block)}
			for _, a := range t.Attrs {
				tag.Attrs = append(tag.Attrs, Attribute{Name: a.Name, Val: JavaScriptExpression(fmt.Sprintf("%v", a.Val)), MustEscape: a.MustEscape})
			}

			res = append(res, tag)

		case "Mixin":
			mixin := new(Mixin)
			mixin.Block = Block{Nodes: p.build(t.Block)}
			for _, a := range t.AttributeBlocks {
				mixin.AttributeBlocks = append(mixin.AttributeBlocks, JavaScriptExpression(a))
			}
			mixin.Name = JavaScriptIdentifier(t.Name)
			mixin.Args = t.Args
			for _, a := range t.Attrs {
				mixin.Attrs = append(mixin.Attrs, Attribute{Name: a.Name, Val: JavaScriptExpression(fmt.Sprintf("%v", a.Val)), MustEscape: a.MustEscape})
			}
			mixin.Call = t.Call

			res = append(res, mixin)

		case "Text":
			text := new(Text)
			text.Val = t.Val
			res = append(res, text)

		case "Code":
			code := new(Code)
			code.Val = t.Val
			code.Block = Block{Nodes: p.build(t.Block)}
			code.IsInline = t.IsInline
			code.MustEscape = t.MustEscape
			res = append(res, code)

		case "Conditional":
			cond := new(Conditional)
			cond.Test = JavaScriptExpression(t.Test)
			cond.Consequent = Block{Nodes: p.build(t.Consequent)}
			if t.Alternate != nil {
				cond.Alternate = Block{Nodes: p.build(t.Alternate)}
			}
			res = append(res, cond)

		case "Each":
			each := new(Each)
			each.Val = JavaScriptIdentifier(t.Val)
			each.Key = JavaScriptIdentifier(t.Key)
			each.Obj = JavaScriptExpression(t.Obj)
			each.Block = Block{Nodes: p.build(t.Block)}

			res = append(res, each)

		case "Doctype":
			doctype := new(Doctype)
			doctype.Val = t.Val

			res = append(res, doctype)

		case "NamedBlock", "Block":
			res = append(res, &Block{Nodes: p.build(t)})

		default:
			log.Printf("%#v\n", t)
			panic(t)
		}
	}
	return
}
