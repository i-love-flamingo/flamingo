package search

import (
	"context"
	"encoding/json"
	categorydomain "flamingo/core/category/domain"
	productdomain "flamingo/core/product/domain"
	"flamingo/core/search/domain"
	productdto "flamingo/om3/searchperience/infrastructure/product/dto"
	"flamingo/om3/searchperience/infrastructure/search/dto"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type (
	// Service for service usage
	Service struct {
		Client  *Client `inject:""`
		Locale  string  `inject:"config:locale"`
		Channel string  `inject:"config:searchperience.frontend.channel"`
	}
)

var (
	_ domain.SearchService = new(Service)
)

// Search a result
func (s *Service) GetProducts(ctx context.Context, searchMeta domain.SearchMeta, filter ...domain.Filter) (domain.SearchMeta, []productdomain.BasicProduct, []domain.Filter, error) {
	var categoryRequest string

	query := url.Values{
		"channel": {s.Channel},
		"locale":  {s.Locale},
	}

	for _, v := range filter {
		if categoryFacet, ok := v.(*categorydomain.CategoryFacet); ok {
			categoryRequest = categoryFacet.Values()[string(categorydomain.CategoryKey)][0]
			continue
		}
		for k, v := range v.Values() {
			for _, vv := range v {
				query.Add(k, vv)
			}
		}
	}

	var resp *http.Response
	var err error

	if categoryRequest != "" {
		resp, err = s.Client.Category(ctx, categoryRequest, query)
	} else {
		resp, err = s.Client.Search(ctx, query)
	}

	if err != nil {
		return searchMeta, nil, nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return searchMeta, nil, nil, domain.ErrNotFound
	}

	res := new(dto.Result)
	err = json.NewDecoder(resp.Body).Decode(res)

	if err != nil {
		return searchMeta, nil, nil, errors.WithStack(err)
	}

	searchMeta.Page = res.Results.Product.PageInfo.CurrentPage
	searchMeta.NumPages = res.Results.Product.PageInfo.TotalPages
	searchMeta.NumResults = len(res.Results.Product.Hits)

	products := make([]productdomain.BasicProduct, searchMeta.NumResults)
	for i, v := range res.Results.Product.Hits {
		products[i], _ = productdto.Map(ctx, v.Document)
	}

	return searchMeta, products, filter, nil
}
