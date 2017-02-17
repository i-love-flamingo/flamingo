package app

// Register flamingo json handler
func Register(r *ServiceContainer) {
	r.Route("/_flamingo/json/{handler}", "_flamingo.json")
	r.Handle("_flamingo.json", new(GetController))

	r.Register(new(GetFunc), "pug-template.func")
	r.Register(new(UrlFunc), "pug-template.func")
}
