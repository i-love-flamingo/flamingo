package product

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"go.aoe.com/flamingo/core/product/domain"

	"github.com/pkg/errors"
)

var (
	brands = []string{
		"Apple",
		"Bose",
		"Dior",
		"Hugo Boss",
	}
)

// FakeProductService is just mocking stuff
type FakeProductService struct{}

// Get returns a product struct
func (ps *FakeProductService) Get(ctx context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	//defer ctx.Profile("service", "get product "+foreignId)()

	if marketplaceCode == "fake_configurable" {
		product := fakeConfigurable(marketplaceCode)
		product.Title = "TypeConfigurable product"

		product.VariantVariationAttributes = []string{"color", "size"}

		variants := []struct {
			marketplacecode string
			title           string
			attributes      domain.Attributes
		}{
			{"shirt-white-s", "Shirt White S", domain.Attributes{"size": "S", "color": "white"}},
			{"shirt-red-s", "Shirt Red S", domain.Attributes{"size": "S", "color": "red"}},
			{"shirt-white-m", "Shirt White M", domain.Attributes{"size": "M", "color": "white"}},
			{"shirt-black-m", "Shirt Black M", domain.Attributes{"size": "M", "color": "black"}},
			{"shirt-black-l", "Shirt Black L", domain.Attributes{"size": "L", "color": "black"}},
			{"shirt-red-l", "Shirt Red L", domain.Attributes{"size": "L", "color": "red"}},
		}

		for _, variant := range variants {
			simpleVariant := fakeVariant(variant.marketplacecode)
			simpleVariant.Title = variant.title
			simpleVariant.Attributes = variant.attributes

			product.Variants = append(product.Variants, simpleVariant)
		}

		return product, nil
	}
	if marketplaceCode == "fake_simple" {
		return FakeSimple(marketplaceCode), nil
	}
	return nil, errors.New("Not implemented in FAKE: Only code 'fake_configurable' or 'fake_simple' should be used")

}

func FakeSimple(marketplaceCode string) domain.SimpleProduct {
	product := domain.SimpleProduct{}
	product.Title = "TypeSimple product"
	addBasicData(&product.BasicProductData)

	product.ActivePrice = getPrice(20.99+float64(rand.Intn(10)), 10.49+float64(rand.Intn(10)))
	product.MarketPlaceCode = marketplaceCode

	return product
}

func fakeConfigurable(marketplaceCode string) domain.ConfigurableProduct {
	product := domain.ConfigurableProduct{}
	product.Title = "TypeSimple product"
	addBasicData(&product.BasicProductData)
	product.MarketPlaceCode = marketplaceCode

	return product
}

func fakeVariant(marketplaceCode string) domain.Variant {
	var simpleVariant domain.Variant
	simpleVariant.Attributes = make(map[string]interface{})

	addBasicData(&simpleVariant.BasicProductData)

	simpleVariant.ActivePrice = getPrice(30.99+float64(rand.Intn(10)), 20.49+float64(rand.Intn(10)))
	simpleVariant.MarketPlaceCode = marketplaceCode

	return simpleVariant
}

func addBasicData(product *domain.BasicProductData) {
	product.ShortDescription = "Short Description"
	product.Description = "Description"
	product.Media = append(product.Media, domain.Media{Type: "image-external", Reference: "http://pipsum.com/1024x768.jpg?" + strconv.Itoa(rand.Intn(10)), Usage: "detail"})
	product.Media = append(product.Media, domain.Media{Type: "image-external", Reference: "http://pipsum.com/1024x768.jpg?" + strconv.Itoa(rand.Intn(10)), Usage: "list"})
	product.Attributes = make(map[string]interface{})
	product.Attributes["brandCode"] = brands[rand.Intn(len(brands))]
	product.RetailerCode = "Testretailer"
	product.RetailerSku = "12345sku"
	product.CategoryPath = []string{"Testproducts", "Testproducts/Fake/Configurable"}
}

func getPrice(defaultP float64, discounted float64) domain.PriceInfo {
	defaultP, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", defaultP), 64)
	discounted, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", discounted), 64)

	var price domain.PriceInfo
	price.Currency = "EUR"
	price.Default = defaultP
	if discounted != 0 {
		price.Discounted = discounted
		price.DiscountText = "Super test campaign"
		price.IsDiscounted = true
	}
	price.ActiveBase = 1
	price.ActiveBaseAmount = 10
	price.ActiveBaseUnit = "ml"
	return price
}
