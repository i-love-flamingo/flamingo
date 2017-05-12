package interfaces

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"flamingo/om3/search/domain"
)

type (
	// ViewController demonstrates a search view controller
	ViewController struct {
		*responder.ErrorAware    `inject:""`
		*responder.RenderAware   `inject:""`
		*responder.RedirectAware `inject:""`
		domain.SearchService     `inject:""`
	}

	// ViewData is used for search rendering
	ViewData struct {
		SearchResult map[string]interface{}
		SearchHost   string
	}
)

func getSearchType(st string) string {
	switch st {
	case
		"retailer",
		"location",
		"brand":
		return st
	}
	return "product"
}

// Get Response for search
func (vc *ViewController) Get(c web.Context) web.Response {
	query, queryErr := c.Query1("q")
	searchType := getSearchType(c.MustParam1("type"))

	if searchType != c.MustParam1("type") {
		return vc.Redirect("search.view?q="+query, "type", searchType)
	}

	vd := ViewData{
		SearchResult: map[string]interface{}{
			"type":  getSearchType(c.MustParam1("type")),
			"query": query,
		},
		SearchHost: c.Request().Host,
	}

	if query == "" || queryErr != nil {
		return vc.Render(c, "pages/search/view", vd)
	}

	searchResult, err := vc.SearchService.Search(c, c.Request().URL.Query())
	if err != nil {
		return vc.Error(c, err)
	}

	vd.SearchResult["results"] = map[string]interface{}{
		"product":  searchResult.Results.Product,
		"brand":    searchResult.Results.Brand,
		"location": searchResult.Results.Location,
		"retailer": searchResult.Results.Retailer,
	}

	// render page
	return vc.Render(c, "pages/search/view", vd)
}
