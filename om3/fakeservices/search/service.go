package search

//go:generate go-bindata -pkg search -prefix mocks/ mocks/

import (
	"context"
	"strconv"

	productdomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/core/search/domain"
	"go.aoe.com/flamingo/om3/fakeservices/product"
)

type (
	// FakeSearchService is just mocking stuff
	FakeSearchService struct{}
)

var (
	_ domain.SearchService = new(FakeSearchService)
)

func (searchservice *FakeSearchService) GetProducts(ctx context.Context, searchMeta domain.SearchMeta, filter ...domain.Filter) (domain.SearchMeta, []productdomain.BasicProduct, []domain.Filter, error) {
	searchMeta.NumResults = 30
	searchMeta.NumPages = 20
	products := make([]productdomain.BasicProduct, 30)
	for i := 0; i < 30; i++ {
		products[i] = product.FakeSimple("product-" + strconv.Itoa(i))
	}
	return searchMeta, products, filter, nil
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
