package app

// Register flamingo json handler
func Register(r *Registrator) {
	r.Route("/_flamingo/json/{handler}", "_flamingo.json")
	r.Handle("_flamingo.json", new(GetController))
}
