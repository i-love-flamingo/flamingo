package pugast

import "html/template"

// PugAst holds information about the pug abstract syntax tree
type PugAst struct {
	Path    string
	TplCode map[string]string
	mixin   map[string]*Token
	FuncMap template.FuncMap
	rawmode bool
	Doctype string
}

// NewPugAst creates a new Pug AST struct
func NewPugAst(path string) *PugAst {
	pugast := &PugAst{
		Path:    path,
		TplCode: make(map[string]string),
		mixin:   make(map[string]*Token),
	}
	return pugast
}
