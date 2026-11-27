if flamingo.os.env.CONTRACT_REDIS_ADDR != _|_ {
	flamingo: session: {
		backend: "redis"
		redis: host: flamingo.os.env.CONTRACT_REDIS_ADDR
	}
}
