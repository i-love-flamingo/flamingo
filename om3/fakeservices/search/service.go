package search

//go:generate go-bindata -pkg search -prefix mocks/ mocks/

import (
	"encoding/json"
	"flamingo/framework/web"
	"flamingo/om3/search/domain"
	"net/url"
	"strconv"
)

// FakeSearchService is just mocking stuff
type FakeSearchService struct{}

func (searchservice *FakeSearchService) Search(ctx web.Context, query url.Values) (*domain.SearchResult, error) {
	var s = new(domain.SearchResult)
	b, _ := Asset("searchResult.mock.json")
	json.Unmarshal(b, s)

	if page := query.Get("page"); page != "" {
		s.Results.Product.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
		s.Results.Brand.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
		s.Results.Location.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
		s.Results.Retailer.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
	}

	return s, nil
}
