package auth

import (
	"bytes"
	"context"
	"html/template"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"
)

type debugController struct {
	responder       *web.Responder
	identityService *WebIdentityService
}

// Inject dependencies
func (c *debugController) Inject(responder *web.Responder, identityService *WebIdentityService) {
	c.responder = responder
	c.identityService = identityService
}

var tpl = template.Must(template.New("debug").Parse(
	//language=gohtml
	`
<h1>Auth Debug</h1><hr/>
<h2>Registered RequestIdentifier:</h2>
<br/>
{{ range .Identifier }}
{{ .Broker }}: {{ . }} <a href="?__debug__action=authenticate&__debug__broker={{ .Broker }}">Authenticate</a> | <a href="?__debug__action=forceauthenticate&__debug__broker={{ .Broker }}">Force Authenticate</a>
<br/>
{{ end }}
<hr/>
<h2>Active Identities</h2>
{{ range .Identities }}
{{ .Broker}}/{{ .Subject }}: {{ printf "%s / %#v" . . }}
<br/>
{{ end }}
<hr/>
`))

// Action handles auth debugging
func (c *debugController) Action(ctx context.Context, request *web.Request) web.Result {
	action, _ := request.Query1("__debug__action")
	switch action {
	case "forceauthenticate":
		broker, _ := request.Query1("__debug__broker")
		return c.identityService.AuthenticateFor(broker, ctx, request)
	case "authenticate":
		broker, _ := request.Query1("__debug__broker")
		if c.identityService.IdentifyFor(broker, ctx, request) != nil {
			break
		}
		return c.identityService.AuthenticateFor(broker, ctx, request)
	}

	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, struct {
		Identifier []RequestIdentifier
		Identities []Identity
	}{
		Identifier: c.identityService.identityProviders,
		Identities: c.identityService.IdentifyAll(ctx, request),
	})
	if err != nil {
		return c.responder.ServerError(err)
	}

	return c.responder.HTTP(http.StatusOK, buf)
}
