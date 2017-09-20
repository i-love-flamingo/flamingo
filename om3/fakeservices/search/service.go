package search

//go:generate go-bindata -pkg search -prefix mocks/ mocks/

import (
	"context"
	productdomain "flamingo/core/product/domain"
	"flamingo/core/search/domain"
	"flamingo/om3/fakeservices/product"
)

type (
	// FakeSearchService is just mocking stuff
	FakeSearchService struct{}
)

var (
	_ domain.SearchService = new(FakeSearchService)
)

func (searchservice *FakeSearchService) GetProducts(ctx context.Context, searchMeta domain.SearchMeta, filter ...domain.Filter) (domain.SearchMeta, []productdomain.BasicProduct, []domain.Filter, error) {
	ps := new(product.FakeProductService)
	p, _ := ps.Get(ctx, "fake_simple")
	searchMeta.NumResults = 5
	searchMeta.NumPages = 20
	return searchMeta, []productdomain.BasicProduct{p, p, p, p, p}, filter, nil
}

//func (searchservice *FakeSearchService) Search(ctx web.Context, query url.Values) (*domain.SearchResult, error) {
//	var s = new(domain.SearchResult)
//	b, _ := Asset("searchResult.mock.json")
//	json.Unmarshal(b, s)
//
//	if page := query.Get("page"); page != "" {
//		s.Results.Product.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
//		s.Results.Brand.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
//		s.Results.Location.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
//		s.Results.Retailer.PageInfo.CurrentPage, _ = strconv.Atoi(query.Get("page"))
//	}
//
//	return s, nil
//}
