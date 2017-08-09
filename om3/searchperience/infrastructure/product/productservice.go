package product

import (
	"context"
	"encoding/json"
	"flamingo/core/product/domain"
	"flamingo/om3/searchperience/infrastructure/product/dto"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// ProductService for service usage
	ProductService struct {
		Client                 *ProductApiClient  `inject:""`
		Locale                 string             `inject:"config:locale"`
		Channel                string             `inject:"config:searchperience.frontend.channel"`
		TempPriceEngineService PriceEngineService `inject:""`
	}
)

// Get a Product
func (ps *ProductService) Get(ctx context.Context, ID string) (domain.BasicProduct, error) {

	ID = fmt.Sprintf("%s%s%s", ID, ps.Locale, ps.Channel)
	resp, err := ps.Client.Get(ctx, ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.WithStack(domain.ProductNotFound{ID: ID})
	}

	productDto := &dto.Product{}
	err = json.NewDecoder(resp.Body).Decode(productDto)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	mapper := Mapper{}
	return mapper.Map(ctx, productDto, ps.TempPriceEngineService)
}
