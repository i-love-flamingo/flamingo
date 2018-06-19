package application

import (
	"net/url"
	"strings"

	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	Service struct {
		Router  *router.Router `inject:""`
		BaseUrl string         `inject:"config:canonicalurl.baseurl"`
	}
)

func (s *Service) GetBaseDomain() string {
	url, err := url.Parse(s.BaseUrl)

	if err != nil {
		panic(err)
	}

	return url.Host
}

func (s *Service) GetBaseUrl() string {
	return strings.TrimRight(s.BaseUrl, "/")
}

// @todo: Add logic to add allowed parameters via controller
func (s *Service) GetCanonicalUrlForCurrentRequest(ctx web.Context) string {
	return s.GetBaseUrl() + s.Router.Base().Path + ctx.Request().URL.Path
}
