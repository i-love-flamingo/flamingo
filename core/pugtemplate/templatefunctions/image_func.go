package templatefunctions

import (
	"html/template"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
)

type (
	// AssetFunc returns the proper URL for the asset, either local or via CDN
	ImageFunc struct {
		BaseUrl string `inject:"config:imageservice.base_url"`
		Secret  string `inject:"config:imageservice.secret"`
	}
)

// Name alias for use in template
func (imgf ImageFunc) Name() string {
	return "image"
}

// Func as implementation of imageservice helper method
func (imgf *ImageFunc) Func() interface{} {
	return func(source, options, image pugjs.String) template.URL {
		validSources := map[string]bool{
			"pim": true,
			"mdp": true,
			"cms": true,
		}
		if !validSources[source.String()] {
			return ""
		}

		resource := options.String() + "/" + image.String()
		signature := createSignature(resource, imgf.Secret)

		return template.URL(imgf.BaseUrl + "/" + source.String() + "/" + signature + "/" + resource)
	}
}

func createSignature(input string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(input))
	want := mac.Sum(nil)

	return base64.URLEncoding.EncodeToString(want)
}
