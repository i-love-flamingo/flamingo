package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"flamingo/framework"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

// ErrNoPact error
var ErrNoPact = errors.New("no pact setup")

// WithPact runs a test with a pact
func WithPact(t *testing.T, target string, fs ...func(*testing.T, *dsl.Pact)) {
	pact, err := pactSetup("flamingo", target)

	if err != nil {
		t.Skip(err)
		return
	}

	for i, f := range fs {
		t.Run("Pact-"+strconv.Itoa(i), func(t *testing.T) { f(t, pact) })
	}

	if err := pactTeardown(pact); err != nil {
		t.Error(err)
	}
}

// pactSetup sets up pact environment for go tests
func pactSetup(consumer, provider string) (*dsl.Pact, error) {
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
		pactdaemonhost = p
	}

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if _, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", pactdaemonhost, pactdaemonport)); err != nil {
		return nil, ErrNoPact
	}

	var pact = &dsl.Pact{
		Port:     pactdaemonport,
		Host:     pactdaemonhost,
		Consumer: consumer,
		Provider: provider,
		LogLevel: "WARN",
	}

	return pact, nil
}

// pactTeardown tears down the pact instance
func pactTeardown(pact *dsl.Pact) error {
	if pactbroker := os.Getenv("PACT_BROKER_HOST"); pactbroker != "" {
		// Write pact to file `<pact-go>/pacts/my_consumer-my_provider.json`
		if err := pact.WritePact(); err != nil {
			return err
		}

		p := dsl.Publisher{}
		file := filepath.Join(pact.PactDir, fmt.Sprintf("%s-%s.json", strings.ToLower(pact.Consumer), strings.ToLower(pact.Provider)))

		err := p.Publish(types.PublishRequest{
			PactURLs:        []string{file},
			PactBroker:      pactbroker,
			ConsumerVersion: framework.VERSION,
			Tags:            []string{strings.ToLower(pact.Consumer), strings.ToLower(pact.Provider)},
			BrokerUsername:  os.Getenv("PACT_BROKER_USERNAME"),
			BrokerPassword:  os.Getenv("PACT_BROKER_PASSWORD"),
		})
		if err != nil {
			return err
		}
	}

	pact.Teardown()
	return nil
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
