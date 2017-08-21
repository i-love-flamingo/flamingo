package product

import (
	"context"
	"flamingo/core/product/domain"

	"github.com/pkg/errors"
)

// FakeProductService is just mocking stuff
type FakeProductService struct{}

// Get returns a product struct
func (ps *FakeProductService) Get(ctx context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	//defer ctx.Profile("service", "get product "+foreignId)()

	if marketplaceCode == "fake_configurable" {
		var product domain.ConfigurableProduct
		product.Title = "TypeConfigurable product"

		addBasicData(&product.BasicProductData)

		//prepare TypeConfigurable
		product.VariantVariationAttributes = append(product.VariantVariationAttributes, "size")

		var simpleVariant domain.Variant
		simpleVariant.Attributes = make(map[string]interface{})
		addSalableData(&simpleVariant.SaleableData)
		addBasicData(&simpleVariant.BasicProductData)

		simpleVariant.Title = "Variant 1 - L"
		simpleVariant.Attributes["size"] = "L"
		simpleVariant.ActivePrice = getPrice(50, 30)
		simpleVariant.MarketPlaceCode = "variant1code"
		product.Variants = append(product.Variants, simpleVariant)

		simpleVariant.Title = "Variant 1 - XL"
		simpleVariant.Attributes["size"] = "XL"
		simpleVariant.ActivePrice = getPrice(60, 0)
		simpleVariant.MarketPlaceCode = "variant2code"
		product.Variants = append(product.Variants, simpleVariant)

		product.MarketPlaceCode = marketplaceCode
		return product, nil
	}
	if marketplaceCode == "fake_simple" {
		product := domain.SimpleProduct{}
		product.Title = "TypeSimple product"
		addBasicData(&product.BasicProductData)
		addSalableData(&product.SaleableData)
		product.ActivePrice = getPrice(20, 10)
		product.MarketPlaceCode = marketplaceCode
		return product, nil
	}
	var product domain.BasicProduct
	return product, errors.New("Not implemented in FAKE: Only code 'fake_configurable' or 'fake_simple' should be used")

}

func addBasicData(product *domain.BasicProductData) {
	product.ShortDescription = "Short Description"
	product.Description = "Description"
	product.Media = append(product.Media, domain.Media{Type: "image-external", Reference: "http://pipsum.com/1024x768.jpg"})
	product.Attributes = make(map[string]interface{})
	product.Attributes["brandCode"] = "Apple"
	product.RetailerCode = "Testretailer"
}

func addSalableData(product *domain.SaleableData) {
	product.RetailerSku = "12345sku"
}

func getPrice(defaultP float64, discounted float64) domain.PriceInfo {
	var price domain.PriceInfo
	price.Currency = "EUR"
	price.Default = defaultP
	if discounted != 0 {
		price.Discounted = discounted
		price.DiscountText = "Super test campaign"
	}
	price.ActiveBase = 1
	price.ActiveBaseAmount = 10
	price.ActiveBaseUnit = "ml"
	return price
}
