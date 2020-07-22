core: auth: web: broker: [core.oauth.legacyAuthIdentifier, core.auth.http & {broker: "web1", realm: "web1"}]

core: oauth: {
	server: "http://127.0.0.1:3351/dex"
	secret: "client1"
	clientid: "client1"
}

flamingo: session: cookie: secure: false
