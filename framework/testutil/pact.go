package testutil

import (
	"encoding/json"
	"flamingo/framework"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

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

// PactEncodeLike helper to encode a struct as a pact like type
func PactEncodeLike(model interface{}) string {
	var data map[string]interface{}

	var tmp, _ = json.Marshal(model)
	json.Unmarshal(tmp, &data)

	var result = make(map[string]json.RawMessage)

	for k, v := range data {
		if reflect.TypeOf(v) != nil && reflect.TypeOf(v).Kind() == reflect.Map {
			v = []byte(PactEncodeLike(v))
		} else {
			v, _ = json.Marshal(v)
		}

		result[k] = json.RawMessage(dsl.Like(string(v.([]byte))))
	}

	tmp, _ = json.MarshalIndent(result, "", "\t")

	return string(tmp)
}

// PactWithInteractions extends the pact's interactions
func PactWithInteractions(pact dsl.Pact, interactions []*dsl.Interaction) dsl.Pact {
	pact.Interactions = append(pact.Interactions, interactions...)

	p := pact
	mockServer := &dsl.MockService{
		BaseURL:  fmt.Sprintf("http://%s:%d", p.Host, p.Server.Port),
		Consumer: p.Consumer,
		Provider: p.Provider,
	}

	for _, interaction := range p.Interactions {
		err := mockServer.AddInteraction(interaction)
		if err != nil {
			panic(err)
		}
	}

	return pact
}
