package testutil

import (
	"os"
	"strconv"

	"fmt"
	"path/filepath"
	"strings"

	"flamingo/framework"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

// PactSetup sets up pact environment for go tests
func PactSetup(consumer, provider string) dsl.Pact {
	var pactdaemonport = 6666
	var pactdaemonhost = "localhost"
	var err error

	if p := os.Getenv("PACT_DAEMON_PORT"); p != "" {
		pactdaemonport, err = strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
	}

	if p := os.Getenv("PACT_DAEMON_HOST"); p != "" {
		pactdaemonport, err = strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
	}

	var pact = dsl.Pact{
		Port:     pactdaemonport,
		Host:     pactdaemonhost,
		Consumer: consumer,
		Provider: provider,
	}

	return pact
}

// PactTeardown tears down the pact instance
func PactTeardown(pact dsl.Pact) {
	// Write pact to file `<pact-go>/pacts/my_consumer-my_provider.json`
	pact.WritePact()

	if pactbroker := os.Getenv("PACT_BROKER_HOST"); pactbroker != "" {
		p := dsl.Publisher{}
		file := filepath.Join(pact.PactDir, fmt.Sprintf("%s-%s.json", strings.ToLower(pact.Consumer), strings.ToLower(pact.Provider)))

		p.Publish(types.PublishRequest{
			PactURLs:        []string{file},
			PactBroker:      pactbroker,
			ConsumerVersion: framework.VERSION,
			Tags:            []string{strings.ToLower(pact.Consumer), strings.ToLower(pact.Provider)},
			BrokerUsername:  os.Getenv("PACT_BROKER_USERNAME"),
			BrokerPassword:  os.Getenv("PACT_BROKER_PASSWORD"),
		})
	}

	pact.Teardown()
}
