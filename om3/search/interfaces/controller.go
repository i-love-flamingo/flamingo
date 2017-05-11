package interfaces

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"flamingo/om3/search/domain"
)

type (
	// ViewController demonstrates a search view controller
	ViewController struct {
		*responder.ErrorAware  `inject:""`
		*responder.RenderAware `inject:""`
		domain.SearchService   `inject:""`
	}

	// ViewData is used for search rendering
	ViewData struct {
		SearchResult map[string]interface{}
	}
)

// Get Response for search
func (vc *ViewController) Get(c web.Context) web.Response {
	searchResult, err := vc.SearchService.Search(c, c.Request().URL.Query())

	// catch error
	if err != nil {
		return vc.Error(c, err)
	}

	// render page
	return vc.Render(c, "pages/search/view", ViewData{SearchResult: map[string]interface{}{
		"type":  c.MustParam1("type"), // @todo: check for valid type
		"query": c.MustQuery1("q"),
		"results": map[string]interface{}{
			"product":  searchResult.Results.Product,
			"brand":    searchResult.Results.Brand,
			"location": searchResult.Results.Location,
			"retailer": searchResult.Results.Retailer,
		},
	}})
}
