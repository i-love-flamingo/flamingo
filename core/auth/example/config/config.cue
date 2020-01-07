core: auth: {
	// web: broker: ["oidc1", "http1", "oidc2", "http2"]
	web: broker: [
		core.auth.oidc & {broker: "oidc1", clientID: "client1", clientSecret: "client1", "endpoint": "http://127.0.0.1:3351/dex"},
		core.auth.http & {broker: "http1", realm: "http1 realm", users: {"user1": "pw1"}},
		core.auth.oidc & {broker: "oidc2", clientID: "client2", clientSecret: "client2", "endpoint": "http://127.0.0.1:3352/dex"},
		core.auth.http & {broker: "http2", realm: "http2 realm", users: {"user2": "pw2"}},
	]
}

flamingo: session: cookie: secure: false
