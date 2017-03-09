package profiler

import (
	"bytes"
	"flamingo/core/flamingo/web"
	"html/template"
	"net/http"
)

type (
	// ProfileController shows information about a requested profile
	ProfileController struct{}
)

const profileTemplate = `<!doctype html>
<html>
<body style="font-family: sans-serif;">
<h2>Profile {{.Context.ID}}</h2>
Time: {{.Duration}}, Start: {{.Start}}
<hr/>
{{ range $entry := .Childs }}
{{ template "entry" $entry }}
{{ end }}
</body>
</html>

{{ define "entry" }}
<div style="padding-left: 30px; border: dashed 1px #eee;">
{{ .Msg }} ({{ .Duration }})<br/>
<span style="font-size: 10pt; color: #888">
{{ .Fnc }}<br/>
{{.File}} Lines: {{ .Startpos }} - {{ .Endpos }}
<div onclick="this.childNodes[1].style.display='block'">
Click to view
<pre style="display:none">
{{ .Filehint }}
</pre>
</div>
</span>
{{ range $entry := .Childs }}
{{ template "entry" $entry }}
{{ end }}
</div>
{{ end }}
`

// Get Response for Debug Info
func (dc *ProfileController) Get(ctx web.Context) web.Response {
	t, err := template.New("tpl").Parse(profileTemplate)
	if err != nil {
		panic(err)
	}
	var body = new(bytes.Buffer)

	t.ExecuteTemplate(body, "tpl", profilestorage[ctx.Param1("Profile")])

	return &web.ContentResponse{
		ContentType: "text/html; charset=utf-8",
		Status:      http.StatusOK,
		Body:        body,
	}
}
