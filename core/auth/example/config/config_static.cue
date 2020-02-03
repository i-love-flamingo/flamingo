core: auth: web: broker: [
	AuthHttp1,
	AuthHttp2,
	StaticAuthBroker & {broker: "static1", users: ["user1", "user2"]},
]
