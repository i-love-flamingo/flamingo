package auth

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/url"

	"flamingo.me/flamingo/v3/framework/web"
)

type debugController struct {
	responder       *web.Responder
	identityService *WebIdentityService
	reverseRouter   web.ReverseRouter
}

// Inject dependencies
func (c *debugController) Inject(responder *web.Responder, identityService *WebIdentityService, reverseRouter web.ReverseRouter) {
	c.responder = responder
	c.identityService = identityService
	c.reverseRouter = reverseRouter
}

var tpl = template.Must(template.New("debug").Parse(
	//language=gohtml
	`
<h1>Auth Debug</h1><hr/>
<h2>Registered RequestIdentifier:</h2>
<br/>
{{ range .Identifier }}
{{ .Broker }}: {{ . }} <a href="?__debug__action=authenticate&__debug__broker={{ .Broker }}">Authenticate</a> | <a href="?__debug__action=forceauthenticate&__debug__broker={{ .Broker }}">Force Authenticate</a><br />
{{ end }}
<hr/>
<h2>Active Identities</h2>
<a href="?__debug__action=logoutall">Logout All</a><br/>
{{ range .Identities }}
{{ if .Identity }}
{{ .Broker }}: {{ .Identity.Broker}}/{{ .Identity.Subject }}: {{ printf "%T: %s" .Identity .Identity }} <a href="?__debug__action=logout&__debug__broker={{ .Broker }}">Logout</a><br />
{{ else }}
{{ .Broker }}: {{ .Error }}<br/>
{{ end }}
{{ end }}
<hr/>
`))

// Action handles auth debugging
func (c *debugController) Action(ctx context.Context, request *web.Request) web.Result {
	u, _ := c.reverseRouter.Absolute(request, request.Handler.GetHandlerName(), nil)
	request.Params["redirecturl"] = u.String()

	action, _ := request.Query1("__debug__action")
	switch action {
	case "forceauthenticate":
		broker, _ := request.Query1("__debug__broker")

		return c.identityService.AuthenticateFor(ctx, broker, request)

	case "authenticate":
		broker, _ := request.Query1("__debug__broker")
		if identity, _ := c.identityService.IdentifyFor(ctx, broker, request); identity != nil {
			break
		}

		return c.identityService.AuthenticateFor(ctx, broker, request)

	case "logoutall":
		return c.identityService.Logout(ctx, request, &url.URL{Path: request.Request().URL.Path, ForceQuery: true})

	case "logout":
		broker, _ := request.Query1("__debug__broker")

		return c.identityService.LogoutFor(ctx, broker, request, &url.URL{Path: request.Request().URL.Path, ForceQuery: true})
	}

	type identityInfo struct {
		Broker   string
		Identity Identity
		Error    error
	}
	identities := make([]identityInfo, len(c.identityService.identityProviders))

	for i, ip := range c.identityService.identityProviders {
		identity, err := ip.Identify(ctx, request)
		identities[i] = identityInfo{Broker: ip.Broker(), Identity: identity, Error: err}
	}

	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, struct {
		Identifier []RequestIdentifier
		Identities []identityInfo
	}{
		Identifier: c.identityService.identityProviders,
		Identities: identities,
	})
	if err != nil {
		return c.responder.ServerError(err)
	}

	return c.responder.HTTP(http.StatusOK, buf)
}
