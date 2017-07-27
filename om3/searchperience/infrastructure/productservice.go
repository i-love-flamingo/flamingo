package infrastructure

import (
	"context"
	"encoding/json"
	"flamingo/core/product/domain"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// ProductService for service usage
	ProductService struct {
		Client  *ProductClient `inject:""`
		Locale  string         `inject:"config:locale"`
		Channel string         `inject:"config:searchperience.frontend.channel"`
	}
)

// Get a Product
func (ps *ProductService) Get(ctx context.Context, ID string) (*domain.Product, error) {
	ID = fmt.Sprintf("%s_%s_%s", ID, ps.Locale, ps.Channel)
	resp, err := ps.Client.Get(ctx, ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("Product not found")
	}

	res := &domain.Product{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}
