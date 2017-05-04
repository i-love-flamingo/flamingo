package infrastructure

import (
	"context"
	"encoding/json"
	"flamingo/framework/web"
	"flamingo/om3/brand/domain"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// BrandService for service usage
	BrandService struct {
		Client *BrandsClient `inject:""`
	}
)

// Get a brand
func (bs *BrandService) Get(ctx context.Context, ID string) (*domain.Brand, error) {
	if ctx, ok := ctx.(web.Context); ok {
		defer ctx.Profile("masterdataportal", "get brand "+ID)()
	}

	resp, err := bs.Client.Get(ctx, ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.WithStack(domain.BrandNotFound{Name: ID})
	}

	res := &domain.Brand{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}
