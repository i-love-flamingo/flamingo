package infrastructure

import (
	"flamingo/framework/testutil"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
)

var pact dsl.Pact

func TestMain(m *testing.M) {
	pact = testutil.PactSetup("Flamingo", "Searchperience")

	status := m.Run()

	testutil.PactTeardown(pact)

	os.Exit(status)
}
