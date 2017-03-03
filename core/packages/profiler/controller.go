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
<body>
<pre>
{{.}}
</pre>
</body>
</html>
`

// Get Response for Debug Info
func (dc *ProfileController) Get(ctx web.Context) web.Response {
	t, _ := template.New("tpl").Parse(profileTemplate)
	var body = new(bytes.Buffer)

	t.Execute(body, profilestorage[ctx.Param1("Profile")].String())

	return &web.ContentResponse{
		ContentType: "text/html",
		Status:      http.StatusOK,
		Body:        body,
	}
}
