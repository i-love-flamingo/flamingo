package specs

import (
	"context"
	"flamingo/core/product/domain"
	"flamingo/framework/testutil"
	"flamingo/om3/searchperience/infrastructure/product"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
)

func TestProductserviceCanGetSimpleProduct(t *testing.T) {
	testutil.WithPact(t, "searchperience-frontend", func(pact *dsl.Pact) {
		var simpleTitleFixture = "TED BAKER Cro Polo"

		pact.AddInteraction().
			Given("The simple test product exists").
			UponReceiving("A request to a simple test product").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   "/product/TEDBAKERCroPolo_simple-en_EN-onlinestore",
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Body:   getSimpleProductResponseFixture(),
			})

		if err := pact.Verify(func() error {

			var productclient = product.ProductApiClient{}
			productclient.SearchperienceClient.BaseURL = fmt.Sprintf("http://localhost:%d/", pact.Server.Port)
			var productService = product.ProductService{Client: &productclient, Locale: "en_EN", Channel: "onlinestore"}

			var testProduct, err = productService.Get(context.Background(), "TEDBAKERCroPolo_simple")

			if err != nil {
				t.Error(err)
			}

			if testProduct == nil {
				t.Error("Product is nil")
			}
			if testProduct.BaseData().Title != simpleTitleFixture {
				t.Errorf("Product Title is expected to be %v got %v", simpleTitleFixture, testProduct.BaseData().Title)
			}
			return nil
		}); err != nil {
			t.Error(err)
		}
	})
}

func getSimpleProductResponseFixture() string {
	b, _ := ioutil.ReadFile("product/fixture.simple.json")
	return string(b)
}

func TestProductserviceCanGetConfigurableProduct(t *testing.T) {
	testutil.WithPact(t, "searchperience-frontend", func(pact *dsl.Pact) {
		var titleFixture = "TUMI Travel Accessories Large Packing Cube"
		pact.AddInteraction().
			Given("The configurable test product exists").
			UponReceiving("A request to a configurable test product").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   "/product/TUMITravelAccessoriesLargePackingCube_configurable-en_EN-mainstore",
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
			})

		if err := pact.Verify(func() error {

			var productclient = product.ProductApiClient{}
			productclient.SearchperienceClient.BaseURL = fmt.Sprintf("http://localhost:%d/", pact.Server.Port)
			var productService = product.ProductService{Client: &productclient, Locale: "en_EN", Channel: "mainstore"}

			var testProduct, err = productService.Get(context.Background(), "TUMITravelAccessoriesLargePackingCube_configurable")

			if err != nil {
				t.Error(err)
			}

			if testProduct == nil {
				t.Errorf("Product is nil")
			}

			//if testProduct.BaseData().Title != titleFixture {
			//	t.Errorf("Product Title is expected to be %v got %v", titleFixture, testProduct.BaseData().Title)
			//}

			configurableProduct := testProduct.(domain.ConfigurableProduct)
			if len(configurableProduct.Variants) != 2 {
				t.Errorf("TypeConfigurable product should have 2 Variants")
			}

			if configurableProduct.Variants[0].Title != titleFixture {
				t.Errorf("Variant Product Title is expected to be %v got %v", titleFixture, configurableProduct.Variants[0].Title)
			}

			// if configurableProduct.Variants[0].ActivePrice.Default != 0 {
			// 	t.Errorf("Variant Product ActivePrice is expected to not be 0 Got: %v", configurableProduct.Variants[0].ActivePrice.Default)
			// }

			return nil
		}); err != nil {
			t.Error(err)
		}
	})
}

func getConfigurableProductResponseFixture() string {
	b, _ := ioutil.ReadFile("product/fixture.configurable.json")
	return string(b)
}
