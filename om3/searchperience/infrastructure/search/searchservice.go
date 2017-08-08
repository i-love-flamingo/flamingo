package infrastructure

import (
	"encoding/json"
	"flamingo/framework/web"
	"flamingo/om3/search/domain"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type (
	// SearchService for service usage
	SearchService struct {
		Client *SearchClient `inject:""`
	}
)

// Search a result
func (ss *SearchService) Search(ctx web.Context, query url.Values) (*domain.SearchResult, error) {
	resp, err := ss.Client.Search(ctx, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("Search not available")
	}

	res := &domain.SearchResult{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}
