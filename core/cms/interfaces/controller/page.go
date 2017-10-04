package controller

import (
	"go.aoe.com/flamingo/core/cms/domain"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// ViewController demonstrates a product view controller
	ViewController struct {
		responder.ErrorAware  `inject:""`
		responder.RenderAware `inject:""`
		PageService           domain.PageService `inject:""`
	}

	// ViewData for rendering
	ViewData struct {
		CmsPage domain.Page
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	var page, err = vc.PageService.Get(c, c.MustParam1("name"))

	if page == nil && err == nil {
		return vc.ErrorNotFound(c, err)
	}

	if err != nil {
		return vc.Error(c, err)
	}

	//fmt.Printf("%+v\n", page)

	//res2B, _ := json.Marshal(page)
	//fmt.Println(string(res2B))

	template, err := c.Param1("template")
	if err != nil {
		template = "cms/cms"
	}

	return vc.Render(c, template, ViewData{CmsPage: *page})
}
