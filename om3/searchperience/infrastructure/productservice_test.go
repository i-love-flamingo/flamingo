package infrastructure

import (
	"context"
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
)

func TestProductservice(t *testing.T) {
	pact.AddInteraction().
		Given("The test product exists").
		UponReceiving("A request to the test product").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   "/product/test",
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   `{"identifier":"foo"}`,
		})

	if err := pact.Verify(func() error {
		var productclient = NewProductClient(&SearchperienceClient{
			BaseURL: fmt.Sprintf("http://localhost:%d/", pact.Server.Port),
		})

		var product, err = productclient.Get(context.Background(), "test")

		if err != nil {
			t.Fatal(err)
		}

		if product == nil {
			t.Fatal("Product is nil")
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
