package pug_template

import (
	"bytes"
	"flamingo/core/flamingo/web"
	"flamingo/core/packages/pug_template/pugast"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type (
	DebugController struct {
		Engine *pugast.PugTemplateEngine `inject:""`
	}
)

const DebugTemplate = `<!doctype html>
<html>
<head>
	<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.9.0/styles/default.min.css">
	<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.9.0/highlight.min.js"></script>
</head>

<body>
<pre><code class="html">{{ . }}</code></pre>

<script>hljs.initHighlightingOnLoad();</script>
</body>
</html>
`

func (dc *DebugController) Get(ctx web.Context) web.Response {
	//if dc.Engine.ast == nil {
	//dc.Engine.loadTemplates()
	//}

	tpl, ok := dc.Engine.Ast.TplCode[ctx.Query1("tpl")]
	if !ok {
		panic("tpl not found")
	}
	t, _ := template.New("tpl").Parse(DebugTemplate)
	var body = new(bytes.Buffer)

	tpls := ""
	for i, l := range strings.Split(tpl, "\n") {
		tpls += fmt.Sprintf("%03d: %s\n", (i + 1), l)
	}

	t.Execute(body, tpls)

	return web.ContentResponse{
		ContentType: "text/html",
		Status:      http.StatusOK,
		Body:        body,
	}
}
