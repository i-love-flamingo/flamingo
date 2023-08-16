# Auth

The Auth module provides access to authentication logic.
It comes with the following features:

* Basic models for "Identity" and a concept how to register new "RequestIdentifiers"
* A "WebModule" that can be used to handle and start login flows (or authentication flows)
* Implementation of:
    * OpenIDConnect Authentication Flow - including the possibility to fake authentications for testing and development purposes
    * Identify requests with OAuth Bearer Token 
    * And a simple HTTP Authentication (BasicAuth) implementation


## Usage Examples

### Use OIDC in your (web) application
A typical usecase is, that you want to authenticate users against a given single sign on system via OICD (Open ID Connect).

Therefore add the `oauth.Module` to the bootstrap of your application. (The module will automatically add the `auth.WebModule` as a dependency:

```go

package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/auth/oauth"
)

func main() {
	flamingo.App([]dingo.Module{
		new(oauth.Module),
	})
}
```

Then configure the Broker that you want to use in your application configuration. For example if you want to [use google as an open id connect provider](https://developers.google.com/identity/openid-connect/openid-connect) you can configure it like this:

```yaml
flamingo.debug.mode: true

core.auth.web.broker:
  - 
    broker: "google"
    clientID: **************
    clientSecret: **********
    enableEndSessionEndpoint: true
    # request offline_access scope to receive refresh
    enableOfflineToken: false
    endpoint: "https://accounts.google.com"
    typ: "oidc"


flamingo.session.cookie.secure: false

```
Then you can trigger the authentication process by calling the route: `http://localhost:3322/core/auth/login/google`
You can find this example in the folder "example/google". 

To access Identity informations you can use the `auth.WebIdentityService` inside your code:

```go 
// Inject dependencies
func (controller *testController) Inject(webIdentityService *auth.WebIdentityService) *testController {
	controller.webIdentityService = webIdentityService
	return controller
}

func (controller *testController) Index(ctx context.Context, req *web.Request) web.Result {
	identity := controller.webIdentityService.Identify(ctx, req)
	
	if identity == nil {
		// not logged in
	} else {
        // logged in
		oidcIdentity, _ := identity.(oauth.OpenIDIdentity)
		_ := oidcIdentity.IDToken()
	}
    ...
}

```



### Use multiple Brokers in one application:

There is also an advanced example in the "example" folder, which comes with a docker compose setup booting up two [dex](https://dexidp.io/) and two [keycloak](https://www.aoe.com/techradar/tools/keycloak.html) and a confiured flamingo application that uses all kind of Authentication Providers - including basic auth and fake implementations for testing.

In debug mode (`core.auth.web.debugController`, default to `flamingo.debug.mode`) there is http://localhost:3322/core/auth/debug for debugging.


## Concept and Design

The module uses the following concepts:

* An *Identity* is the object that represents "someone" or "something". It just has a "subject" (e.g. a username or id). An  *Identity* is "identified" by a *Broker*
* In the module we use the following differentiation:
    * Identify: Is the process of checking the request for an existing Identity information. (Some broker may use a state to detect an Identity)
    * Authenticate: Is the process of requesting an Identity. (e.g. by providing a login form or redirection to an external Authprovider)
* A *Broker* is a specific implementation of an Authorisation and Identification logic/sheme. A Broker may just implement the basic *RequestIdentifier* interface or also *WebAuthenticater*, *WebCallbacker*, *WebLogouter* and *WebIdenityRefresher* interfaces if the implemented Authorisationlogic needs it. 
* Identity Brokers have a type and can be registered and configured under a "Brokername". In general, all identifying brokers types can be configured more than once, and at any point, there can be zero to many identities available.

More details in the following chapters:

### Identity
The `auth.Identity` is the minimum available information about an identified request/context situation.
It consists of `Broker` and `Subject`, where the broker identifies the authenticating broker (authenticating party) and the subject identifies the primary subject the identity identifies.

### RequestIdentifier
The WebIdentifier primarily identifies incoming `web.Request`s.
This could be done by means of inspecting the session, request data (auth header), etc.

### WebAuthenticator
WebIdentifier who implements the authenticator interface is able to trigger authentication. This can be a redirect to an external page, setting HTTP headers, or presenting a login form.

### WebLogouter
Once a logout has triggered all identifiers who implement either one of the logout methods are called.
The WebLogouter will destroy session data etc., while the WebLogouterWithRedirect can return a redirect (e.g. to an OpenID Connect server).

Multiple redirects are handled automagically.


