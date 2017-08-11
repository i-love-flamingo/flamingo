package specs

import (
	"context"
	"testing"

	product "flamingo/om3/searchperience/infrastructure/product"
	"fmt"

	"io/ioutil"

	"flamingo/core/product/domain"

	"github.com/pact-foundation/pact-go/dsl"
)

func TestProductserviceCanGetSimpleProduct(t *testing.T) {
	var simpleTitleFixture = "simple-title"

	pact.AddInteraction().
		Given("The simple test product exists").
		UponReceiving("A request to a simple test product").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   "/product/testen_USecommerce",
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   getSimpleProductResponseFixture(),
		})

	// @TODO: Remove when serach is ready: This is just for the time where search is not calling the priceengine
	pact.AddInteraction().
		Given("TEMP: Priceegine product exists").
		UponReceiving("A request to a priceengine for simple").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   "/prices",
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   getPriceResponseFixture("searchperience-product-v2"),
		})

	if err := pact.Verify(func() error {

		var productclient = product.ProductApiClient{}
		productclient.SearchperienceClient.BaseURL = fmt.Sprintf("http://localhost:%d/", pact.Server.Port)
		var productService = product.ProductService{Client: &productclient, Locale: "en_US", Channel: "ecommerce"}
		productService.TempPriceEngineService.TempPriceEngineBaseURL = fmt.Sprintf("http://localhost:%d/", pact.Server.Port)

		var product, err = productService.Get(context.Background(), "test")

		if err != nil {
			t.Fatal(err)
		}

		if product == nil {
			t.Fatal("Product is nil")
		}
		if product.GetBaseData().Title != simpleTitleFixture {
			t.Fatalf("Product Title is expected to be %v got %v", simpleTitleFixture, product.GetBaseData().Title)
		}
		// @TODO: Remove when search is ready: This is just for the time where search is not calling the priceengine
		simpleProduct := product.(domain.SimpleProduct)
		if simpleProduct.ActivePrice.Default != 9.99 {
			t.Fatalf("Product ActivePrice is expected to be %v got %v", 9.99, simpleProduct.ActivePrice.Default)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func getSimpleProductResponseFixture() string {
	byte, _ := ioutil.ReadFile("product/fixture.simple.json")
	return string(byte)
}

func TestProductserviceCanGetConfigurableProduct(t *testing.T) {
	var titleFixture = "Bombay Sapphire Gin Configurable"
	pact.AddInteraction().
		Given("The configurable test product exists").
		UponReceiving("A request to a configurable test product").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   "/product/configurabletestiden_USecommerce",
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   getConfigurableProductResponseFixture(),
		})

	// @TODO: Remove when serach is ready: This is just for the time where search is not calling the priceengine
	pact.AddInteraction().
		Given("TEMP: Priceengine product exists").
		UponReceiving("A request to a priceengine for configurable").
		WithRequest(dsl.Request{
			Method: "POST",
			Path:   "/prices",
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   getPriceResponseFixture("searchperience-product-v2"),
		})

	if err := pact.Verify(func() error {

		var productclient = product.ProductApiClient{}
		productclient.SearchperienceClient.BaseURL = fmt.Sprintf("http://localhost:%d/", pact.Server.Port)
		var productService = product.ProductService{Client: &productclient, Locale: "en_US", Channel: "ecommerce"}

		productService.TempPriceEngineService.TempPriceEngineBaseURL = fmt.Sprintf("http://localhost:%d/", pact.Server.Port)

		var product, err = productService.Get(context.Background(), "configurabletestid")

		if err != nil {
			t.Fatal(err)
		}

		if product == nil {
			t.Fatal("Product is nil")
		}
		if product.GetBaseData().Title != titleFixture {
			t.Fatalf("Product Title is expected to be %v got %v", titleFixture, product.GetBaseData().Title)
		}

		configurableProduct := product.(domain.ConfigurableProduct)
		if len(configurableProduct.Variants) != 3 {
			t.Fatalf("Configurable product should have 3 Variants")
		}

		if configurableProduct.Variants[0].Title != "Bombay Sapphire Gin 0.5L" {
			t.Fatalf("Variant Product Title is expected to be %v got %v", "Bombay Sapphire Gin 0.5L", configurableProduct.Variants[0].Title)
		}

		if configurableProduct.Variants[0].ActivePrice.Default == 0 {
			t.Fatalf("Variant Product ActivePrice is expected to not be 0 Got: %v", configurableProduct.Variants[0].ActivePrice.Default)
		}

		//Check result to priceeginge (
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func getConfigurableProductResponseFixture() string {
	byte, _ := ioutil.ReadFile("product/fixture.configurable.json")
	return string(byte)
}

func getPriceResponseFixture(code string) string {
	return `[
	{
		"marketplaceCode": "` + code + `",
		"activePrice": {
	"default": 9.99,
	"discounted": 4.99,
	"discountText": "Five dollars off",
	"activeBase": 4.50,
	"activeBaseAmount": 0,
	"activeBaseUnit": "",
	"context": {
	  "customerGroup": null,
	  "channelCode": "mainstore",
	  "locale": "de_DE"
	}
	},
		"availablePrices": null
	}]`
}
