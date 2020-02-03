package custom

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

// StaticModule configures a custom static broker implementation example
type StaticModule struct{}

// Configure dependency injection
func (*StaticModule) Configure(injector *dingo.Injector) {
	injector.BindMap(new(auth.RequestIdentifierFactory), "customStaticBroker").ToInstance(func(config config.Map) (auth.RequestIdentifier, error) {
		var cfg struct {
			Broker string   `json:"broker"`
			Users  []string `json:"users"`
		}

		if err := config.MapInto(&cfg); err != nil {
			return nil, err
		}

		return &staticAuthBroker{
			broker: cfg.Broker,
			users:  cfg.Users,
		}, nil
	})
}

// CueConfig schema
func (*StaticModule) CueConfig() string {
	return `
StaticAuthBroker :: {
	broker: string
	typ: "customStaticBroker"
	users: [...string]
}
`
}

type staticAuthBroker struct {
	broker        string
	users         []string
	responder     *web.Responder
	reverseRouter web.ReverseRouter
}

func (b *staticAuthBroker) Inject(responder *web.Responder, reverseRouter web.ReverseRouter) {
	b.responder = responder
	b.reverseRouter = reverseRouter
}

func (b *staticAuthBroker) key(suffix string) string {
	return "custom." + b.broker + "." + suffix
}

func (b *staticAuthBroker) Broker() string {
	return b.broker
}

func (b *staticAuthBroker) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	identity, ok := request.Session().Load(b.key("identity"))
	if !ok {
		return nil, errors.New("no identity stored")
	}
	return identity.(*customIdentity), nil
}

func (b *staticAuthBroker) Authenticate(ctx context.Context, request *web.Request) web.Result {
	body := `<h1>Login</h1><hr/>`
	for _, user := range b.users {
		href, err := b.reverseRouter.Relative("core.auth.callback", map[string]string{"broker": b.broker})
		if err != nil {
			return b.responder.ServerError(err)
		}
		query := href.Query()
		query.Set("user", user)
		href.RawQuery = query.Encode()
		body += fmt.Sprintf(`<a href="%s">%s</a><br/>`, href.String(), user)
	}
	return b.responder.HTTP(http.StatusOK, strings.NewReader(body))
}

func (b *staticAuthBroker) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {
	user, err := request.Query1("user")
	if err != nil || user == "" {
		return b.responder.ServerError(errors.New("no user set"))
	}

	request.Session().Store(b.key("identity"), &customIdentity{BrokerKey: b.broker, UserKey: user})

	return b.responder.URLRedirect(returnTo(request))
}

func (b *staticAuthBroker) Logout(ctx context.Context, request *web.Request) {
	request.Session().Delete(b.key("identity"))
}

func init() {
	gob.Register(new(customIdentity))
}

type customIdentity struct {
	BrokerKey string
	UserKey   string
}

func (i *customIdentity) Subject() string {
	return i.UserKey
}

func (i *customIdentity) Broker() string {
	return i.BrokerKey
}
