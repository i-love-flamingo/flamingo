package pugtemplate

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/framework/web"
)

type (
	// DebugController shows the intermediate go-template compiled from pug AST
	DebugController struct {
		Engine *pugjs.Engine `inject:""`
	}
)

const debugTemplate = `<!doctype html>
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

// Get Response for Debug Info
func (dc *DebugController) Get(ctx context.Context, r *web.Request) web.Response {
	dc.Engine.LoadTemplates(r.MustQuery1("tpl"))

	tpl, ok := dc.Engine.TemplateCode[r.MustQuery1("tpl")]
	if !ok {
		panic("tpl not found")
	}
	t, _ := template.New("tpl").Parse(debugTemplate)
	var body = new(bytes.Buffer)

	tpls := ""
	for i, l := range strings.Split(tpl, "\n") {
		tpls += fmt.Sprintf("%03d: %s\n", i+1, strings.TrimSpace(strings.TrimSuffix(l, `{{- "" -}}`)))
	}

	t.Execute(body, tpls)

	return &web.ContentResponse{
		ContentType: "text/html; charset=utf-8",
		BasicResponse: web.BasicResponse{
			Status: http.StatusOK,
		},
		Body: body,
	}
}
