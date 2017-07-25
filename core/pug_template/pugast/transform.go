package pugast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type (
	// Attr is a simple key-value pair
	Attr struct {
		Name       string
		Val        interface{}
		MustEscape bool
	}

	// Fileref is used by include/extends
	Fileref struct {
		Type, Path string
		Line       int
	}

	// Token defines the basic token read by the tokenizer
	// Tokens form a tree, where the beginning root node starts the document
	Token struct {
		// default
		Type, Name string
		Mode, Val  string
		Line       int

		// subblock
		Block *Token
		// subblock childs
		Nodes []*Token

		// specific information
		AttributeBlocks []string
		Attrs           []*Attr
		MustEscape      bool
		File            *Fileref
		Filename        string
		SelfClosing     bool
		IsInline        *bool
		Obj             string
		Key             string

		// mixin
		Call bool   // mixin call?
		Args string // call args

		// if
		Test                  string // if
		Consequent, Alternate *Token // if result

		// Interpolated
		Expr string
	}
)

// Parse parses a filename into a Token-tree
func (p *PugAst) Parse(file string) *Token {
	bytes, err := ioutil.ReadFile(path.Join(p.Path, file) + ".ast.json")

	if err != nil {
		fmt.Println(file)
		panic(err)
	}

	return p.ParseJSON(bytes, file)
}

// ParseJSON parses a json into a Token-tree
func (p *PugAst) ParseJSON(bytes []byte, file string) *Token {
	token := new(Token)

	err := json.Unmarshal(bytes, token)
	if err != nil {
		fmt.Println(file)
		panic(err)
	}

	return token
}

// TokenToTemplate gets named Template from Token
func (p *PugAst) TokenToTemplate(name string, t *Token) *Template {
	tpl := New(name).
		Funcs(funcmap).
		Funcs(p.FuncMap).
		Option("missingkey=error")

	tr := p.build(t)
	tc := ""

	for _, b := range tr {
		bla, _ := b.Render(p, 0)
		tc += bla
	}

	for _, b := range p.mixinblocks {
		tc += "\n" + b
	}

	for _, b := range p.mixin {
		tc += "\n" + b
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
		node := p.buildNode(t)
		if node != nil {
			res = append(res, node)
		}
	}
	return
}

func (p *PugAst) buildNode(t *Token) (res Node) {
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

		return tag

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

		return mixin

	case "Text":
		text := new(Text)
		text.Val = t.Val
		return text

	case "Code":
		code := new(Code)
		code.Val = t.Val
		code.Block = Block{Nodes: p.build(t.Block)}
		code.IsInline = t.IsInline
		code.MustEscape = t.MustEscape
		return code

	case "Conditional":
		cond := new(Conditional)
		cond.Test = JavaScriptExpression(t.Test)
		cond.Consequent = p.buildNode(t.Consequent)
		if t.Alternate != nil {
			cond.Alternate = p.buildNode(t.Alternate)
		}
		return cond

	case "Each":
		each := new(Each)
		each.Val = JavaScriptIdentifier(t.Val)
		each.Key = JavaScriptIdentifier(t.Key)
		each.Obj = JavaScriptExpression(t.Obj)
		each.Block = Block{Nodes: p.build(t.Block)}

		return each

	case "Doctype":
		doctype := new(Doctype)
		doctype.Val = t.Val

		return doctype

	case "NamedBlock", "Block":
		return &Block{Nodes: p.build(t)}

	case "Comment":
		return nil

	case "Case":
		// &pugast.Token{Type:"Case", Name:"", Mode:"", Val:"", Line:64, Block:(*pugast.Token)(0xc420187b80), Nodes:[]*pugast.Token(nil), AttributeBlocks:[]string(nil), Attrs:[]*pugast.Attr(nil), MustEscape:false, File:(*pugast.Fileref)(nil), Filename:"pages/search/view.pug", SelfClosing:false, IsInline:(*bool)(nil), Obj:"", Key:"", Call:false, Args:"", Test:"", Consequent:(*pugast.Token)(nil), Alternate:(*pugast.Token)(nil), Expr:"SearchResult.type"}
		cas := new(Case)
		cas.Expr = JavaScriptExpression(t.Expr)
		cas.Block = Block{Nodes: p.build(t.Block)}

		return cas

	case "When":
		// &pugast.Token{Type:"When", Name:"", Mode:"", Val:"", Line:65, Block:(*pugast.Token)(0xc4206bb400), Nodes:[]*pugast.Token(nil), AttributeBlocks:[]string(nil), Attrs:[]*pugast.Attr(nil), MustEscape:false, File:(*pugast.Fileref)(nil), Filename:"pages/search/view.pug", SelfClosing:false, IsInline:(*bool)(nil), Obj:"", Key:"", Call:false, Args:"", Test:"", Consequent:(*pugast.Token)(nil), Alternate:(*pugast.Token)(nil), Expr:"\"product\""}
		when := new(When)
		when.Expr = JavaScriptExpression(t.Expr)
		when.Block = Block{Nodes: p.build(t.Block)}

		return when

	case "MixinBlock":
		return new(MixinBlock)

	default:
		log.Printf("%#v\n", t)
		panic(t)
	}
}
