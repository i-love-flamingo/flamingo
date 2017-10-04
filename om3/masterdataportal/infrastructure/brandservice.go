package infrastructure

import (
	"context"
	"encoding/json"
	"net/http"

	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/om3/brand/domain"

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
