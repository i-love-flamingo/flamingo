# Auth

# Design

Identifier
- identifies an incoming request
    - http auth header
    - http session / oidc

Authenticator
    - identifier sub interface
    - triggers an authentication (redirect, http 401+www-authenticate, etc)

Identity
- identifies a user
    - subject: user ID
    - source: identifier who generated this identity
