package pugast

type PugAst struct {
	Path    string
	TplCode map[string]string
	mixin   map[string]*Token
}

func NewPugAst(path string) *PugAst {
	return &PugAst{
		Path:    path,
		TplCode: make(map[string]string),
		mixin:   make(map[string]*Token),
	}
}
