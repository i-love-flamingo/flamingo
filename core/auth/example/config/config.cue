AuthHttp1 :: core.auth.http & {broker: "http1", realm: "http1 realm", users: {"user1": "pw1"}}
AuthHttp2 :: core.auth.http & {broker: "http2", realm: "http2 realm", users: {"user2": "pw2"}}

core: auth: web: broker: [
	core.auth.oidc & {broker: "oidc1", clientID: "client1", clientSecret: "client1", "endpoint": "http://127.0.0.1:3351/dex"},
	AuthHttp1,
	core.auth.oidc & {broker: "oidc2", clientID: "client2", clientSecret: "client2", "endpoint": "http://127.0.0.1:3352/dex"},
	AuthHttp2,
	core.auth.oidc & {broker: "kc1", clientID: "client1", clientSecret: "", "endpoint": "http://127.0.0.1:3353/auth/realms/Realm1", enableOfflineToken: false},
	core.auth.oidc & {broker: "kc2", clientID: "client2", clientSecret: "", "endpoint": "http://127.0.0.1:3354/auth/realms/Realm2", enableOfflineToken: false},
	customOidcBroker,
	StaticAuthBroker & {broker: "static1", users: ["user1", "user2"]},
]

flamingo: session: cookie: secure: false
