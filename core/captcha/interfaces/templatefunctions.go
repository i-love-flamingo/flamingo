package interfaces

import (
	"fmt"

	"flamingo.me/flamingo/core/captcha/application"
	"flamingo.me/flamingo/core/captcha/domain"
)

type (
	// CaptchaFunc is a template helper to generate a new captcha
	CaptchaFunc struct {
		len       int
		generator *application.Generator
	}

	// CaptchaImgFunc is a template helper to get the image URL for a captcha
	CaptchaImgFunc struct{}

	// CaptchaSoundFunc is a template helper to get the audio URL for a captcha
	CaptchaSoundFunc struct{}
)

// Inject dependencies
func (f *CaptchaFunc) Inject(
	config *struct {
		Len float64 `inject:"config:captcha.len"`
	},
	g *application.Generator,
) {
	f.len = int(config.Len)
	f.generator = g
}

// Func returns the template function to generate a new captcha
func (f *CaptchaFunc) Func() interface{} {
	return func() *domain.Captcha {
		return f.generator.NewCaptcha(f.len)
	}
}

// Func returns the template function to get the image URL from a given captcha
func (f *CaptchaImgFunc) Func() interface{} {
	return urlFuncFactory("png")
}

// Func returns the template function to get the audio URL from a given captcha
func (f *CaptchaSoundFunc) Func() interface{} {
	return urlFuncFactory("wav")
}

func urlFuncFactory(ext string) func(c *domain.Captcha, options ...bool) string {
	return func(c *domain.Captcha, options ...bool) string {
		prefix := "/captcha/"
		if len(options) > 0 && options[0] {
			prefix = prefix + "download/"
		}
		return fmt.Sprintf("%s%s.%s", prefix, c.Hash, ext)
	}
}
