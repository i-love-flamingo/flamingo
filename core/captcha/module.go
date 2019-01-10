package captcha

import (
	"crypto/rand"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/core/captcha/application"
	"flamingo.me/flamingo/core/captcha/interfaces"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/template"
	"github.com/dchest/captcha"
)

type (
	// Module basic struct
	Module struct{}

	routes struct {
		captchaWidth  int
		captchaHeight int
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	router.Bind(injector, new(routes))

	// we need Generator to be a singleton because the key is stored inside the instance
	injector.Bind((*application.Generator)(nil)).In(dingo.Singleton)
	template.BindFunc(injector, "captcha", new(interfaces.CaptchaFunc))
	template.BindFunc(injector, "captchaImage", new(interfaces.CaptchaImgFunc))
	template.BindFunc(injector, "captchaSound", new(interfaces.CaptchaSoundFunc))

	captcha.SetCustomStore(injector.GetInstance((*application.PseudoStore)(nil)).(*application.PseudoStore))
}

// DefaultConfig for the module
func (m *Module) DefaultConfig() config.Map {
	var key [32]byte

	_, err := rand.Read(key[:])
	if err != nil {
		panic(err)
	}

	return config.Map{
		"captcha.len":                  float64(captcha.DefaultLen),
		"captcha.image.width":          float64(captcha.StdWidth),
		"captcha.image.height":         float64(captcha.StdHeight),
		"captcha.encryptionPassPhrase": string(key[:]),
	}
}

// Inject dependencies
func (r *routes) Inject(
	config *struct {
		Width  float64 `inject:"config:captcha.image.width"`
		Height float64 `inject:"config:captcha.image.height"`
	},
) {
	r.captchaWidth = int(config.Width)
	r.captchaHeight = int(config.Height)
}

// Routes registered by this module
func (r *routes) Routes(registry *router.Registry) {
	registry.Route("/captcha/*n", "_captcha")
	registry.HandleAny("_captcha", router.HTTPAction(captcha.Server(r.captchaWidth, r.captchaHeight)))
}
