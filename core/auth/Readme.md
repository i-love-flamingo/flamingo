# Auth

## Design

In general all identifying broker are able to be specified more than once, and at any point there can be zero to many identities available.

### Identity
The `auth.Identity` is the minimum available information about an identified request/context situation.
It consists of `Broker` and `Subject`, where the broker identifies the authenticating party and the subject identifies the primary subject the identity identifies.

### WebIdentifier
The WebIdentifier primarily identifies incoming `web.Request`s.
This could be done by means of inspecting the session, request data (auth header), etc.

### WebAuthenticator
WebIdentifier who implement the authenticator interface are able to trigger an authentication.
This can be a redirect to an external page, setting http headers or presenting a login form.

### WebLogouter
Once a logout is triggered all identifier who implement either one of the logout methods are called.
The WebLogouter will destroy session data etc., while the WebLogouterWithRedirect can return a redirect (e.g. to an oidc server).

Multiple redirects are handled automagically.

## Debug

In debug mode (`core.auth.web.debugController`, default to `flamingo.debug.mode`) there is http://localhost:3322/core/auth/debug for debugging.
