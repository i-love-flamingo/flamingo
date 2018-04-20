package application

import (
	"strings"

	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	Service struct {
		Router  *router.Router `inject:""`
		BaseUrl string         `inject:"config:canonicalurl.baseurl"`
	}
)

func (s *Service) GetBaseUrl() string {
	return strings.TrimRight(s.BaseUrl, "/")
}

// @todo: Add logic to add allowed parameters via controller
func (s *Service) GetCanonicalUrlForCurrentRequest(ctx web.Context) string {
	return s.GetBaseUrl() + s.Router.Base().Path + ctx.Request().URL.Path
}
