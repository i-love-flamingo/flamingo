package framework

import (
	"bytes"
	flamingotemplate "flamingo/framework/template"
	"flamingo/framework/web"
	"html/template"
	"io"
)

// validate interface
var _ flamingotemplate.Engine = new(Engine)

// Engine for template rendering
type Engine struct {
	Glob string `inject:"config:go_template.glob"`
}

// Render a template
func (e *Engine) Render(context web.Context, name string, data interface{}) io.Reader {
	tpl, _ := template.ParseGlob(e.Glob)

	var res = new(bytes.Buffer)
	tpl.ExecuteTemplate(res, name+".html", data)

	return res
}
