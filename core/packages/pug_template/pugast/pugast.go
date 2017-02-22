package pugast

import "html/template"

type PugAst struct {
	Path     string
	TplCode  map[string]string
	mixin    map[string]*Token
	FuncMap  template.FuncMap
	knownVar map[string]bool
}

func NewPugAst(path string) *PugAst {
	pugast := &PugAst{
		Path:     path,
		TplCode:  make(map[string]string),
		mixin:    make(map[string]*Token),
		knownVar: make(map[string]bool),
	}

	pugast.knownVar["attributes"] = true

	return pugast
}
