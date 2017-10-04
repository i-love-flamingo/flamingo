package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/om3/searchperience/infrastructure/product/dto"

	"github.com/pkg/errors"
)

type (
	// ProductService for service usage
	ProductService struct {
		Client  *ProductApiClient `inject:""`
		Locale  string            `inject:"config:locale"`
		Channel string            `inject:"config:searchperience.frontend.channel"`
	}
)

// Get a Product
func (ps *ProductService) Get(ctx context.Context, ID string) (domain.BasicProduct, error) {

	ID = fmt.Sprintf("%s-%s-%s", ID, ps.Locale, ps.Channel)
	resp, err := ps.Client.Get(ctx, ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.WithStack(domain.ProductNotFound{MarketplaceCode: ID})
	}

	productDto := &dto.Product{}
	err = json.NewDecoder(resp.Body).Decode(productDto)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return dto.Map(ctx, productDto)
}
