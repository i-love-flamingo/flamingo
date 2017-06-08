package infrastructure

import (
	"context"
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
)

func TestProductservice(t *testing.T) {
	pact.AddInteraction().
		Given("A Product exists").
		UponReceiving("A request to a product").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   "/document",
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
