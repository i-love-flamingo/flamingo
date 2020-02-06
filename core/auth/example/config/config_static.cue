core: auth: web: broker: [
	AuthHttp1,
	AuthHttp2,
	StaticAuthBroker & {broker: "static1", users: ["user1", "user2"]},
	AuthFake1,
	AuthFake1WithoutPasswords,
	AuthFake1WithDefaultTemplate,
	AuthFake1WithCustomTemplate
]
