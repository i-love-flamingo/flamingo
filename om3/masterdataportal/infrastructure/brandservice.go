package infrastructure

import (
	"context"
	"encoding/json"
	"flamingo/framework/web"
	"flamingo/om3/brand/domain"
	"fmt"
)

type (
	// BrandService for service usage
	BrandService struct {
		Client *BrandsClient `inject:""`
	}
)

// Get a brand
func (bs *BrandService) Get(ctx context.Context, ID string) *domain.Brand {
	if ctx, ok := ctx.(web.Context); ok {
		defer ctx.Profile("masterdataportal", "get brand "+ID)()
	}

	resp := bs.Client.Get(ctx, ID)
	fmt.Println(resp.Header)
	res := &domain.Brand{}
	json.NewDecoder(resp.Body).Decode(res)
	return res
}
