core: auth: {
	fake :: {
		UserConfig :: {
			password?: string
		}

		typ: "fake"
		broker: string
		loginTemplate?: string
		userConfig: {
			[string]: UserConfig
		}

		validatePassword: bool | *true
		usernameFieldId: string | *"username"
		passwordFieldId: string | *"password"
	}
}
customOidcBroker :: {
	typ: "customOidcBroker"
	broker: "customOidcBroker"
	oidc: {
		core.auth.oidc
		broker: "customOidcBroker"
		clientID: "customOidcBroker"
		clientSecret: "customOidcBroker"
		"endpoint": "http://127.0.0.1:3351/dex"
	}
}
StaticAuthBroker :: {
	broker: string
	typ: "customStaticBroker"
	users: [...string]
}
core: auth: {
	http :: {
		typ: "http"
		broker: string
		realm: string
		users: [string]: string
	}
	oauth2Config :: {
        broker: string
        clientID: string
        clientSecret: string
        endpoint: string
        scopes: [...string] | *["profile", "email"]
        claims: {
            accessToken: { [string]: string }
        }
    }

    oidc :: {
        oauth2Config
        typ: "oidc"
        enableOfflineToken: bool | *true
        claims: {
            idToken: { [string]: string } & {
                sub: string | *"sub"
                email: string | *"email"
                name: string | *"name"
            }
        }
        requestClaims: {
            idToken: [...string]
            userInfo: [...string]
        }
        enableEndSessionEndpoint: bool | *true
    }
}
