package infrastructure

import (
	"context"
	"flamingo/framework/testutil"
	"flamingo/om3/brand/domain"
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
)

func TestBrandservice(t *testing.T) {
	testutil.WithPact(t, func(pact dsl.Pact) {
		var identifier = "test-toblerone"
		var expected = domain.Brand{
			ID: identifier,
		}

		pact.AddInteraction().
			Given("A Braqnd exists").
			UponReceiving("A request to a brand").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   "/brands/" + identifier,
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Body:   testutil.PactEncodeLike(expected),
			})

		if err := pact.Verify(func() error {
			var brandsClient = NewBrandsClient(&APIClient{
				BaseURL: fmt.Sprintf("http://localhost:%d/", pact.Server.Port),
			})

			var brand, err = brandsClient.Get(context.Background(), identifier)

			if err != nil {
				t.Fatal(err)
			}

			if brand == nil {
				t.Fatal("Brand is nil")
			}

			return nil
		}); err != nil {
			t.Fatal(err)
		}
	})
}
