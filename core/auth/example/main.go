package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/http"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/core/requestlogger"
)

func main() {
	flamingo.App([]dingo.Module{
		new(requestlogger.Module),
		new(auth.WebModule),
		new(oauth.Module),
		new(http.Module),
	})
}
