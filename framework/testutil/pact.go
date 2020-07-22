package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
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
func WithPact(t *testing.T, from, to string, fs ...func(*testing.T, *dsl.Pact)) {
	if from == "" {
		from = "flamingo"
	}

	pact, err := pactSetup(from, to)
	if err != nil {
		t.Skip(err)
		return
	}
	// defer the pact teardown
	defer func() {
		if err := pactTeardown(pact); err != nil {
			t.Error(err)
		}
	}()

	for i, f := range fs {
		t.Run("Pact-"+strconv.Itoa(i), func(t *testing.T) { f(t, pact) })
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
	defer pact.Teardown()
	if pactbroker := os.Getenv("PACT_BROKER_HOST"); pactbroker != "" {
		// Write pact to file `<pact-go>/pacts/my_consumer-my_provider.json`
		if err := pact.WritePact(); err != nil {
			return err
		}

		p := dsl.Publisher{}
		file := filepath.Join(pact.PactDir, fmt.Sprintf("%s-%s.json", strings.ToLower(pact.Consumer), strings.ToLower(pact.Provider)))

		err := p.Publish(types.PublishRequest{
			PactURLs:        []string{file},
			PactBroker:      strings.TrimSuffix(pactbroker, "/"),
			ConsumerVersion: os.Getenv("PACT_VERSION"),
			Tags:            append([]string{strings.ToLower(pact.Consumer), strings.ToLower(pact.Provider)}, strings.Split(os.Getenv("PACT_TAGS"), ",")...),
			BrokerUsername:  os.Getenv("PACT_BROKER_USERNAME"),
			BrokerPassword:  os.Getenv("PACT_BROKER_PASSWORD"),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// PactEncodeLike encodes a byte slice from json.Marshal or jsonpb into a pact type-like representation
func PactEncodeLike(model []byte) string {
	var data interface{}
	json.Unmarshal(model, &data)

	return string(pactEncode(data))
}

func pactEncode(data interface{}) json.RawMessage {
	switch data := data.(type) {
	case string:
		data = `"` + data + `"`
		return json.RawMessage(dsl.Like(data))

	case int, float32, float64, bool, uint:
		return json.RawMessage(dsl.Like(data))

	case map[string]interface{}:
		for k, v := range data {
			data[k] = pactEncode(v)
		}
		b, _ := json.Marshal(data)
		return json.RawMessage(dsl.Like(string(b)))

	case []interface{}:
		if len(data) < 1 {
			return json.RawMessage(dsl.EachLike(`null`, 0))
		}
		b, _ := json.Marshal(pactEncode(data[0]))
		return json.RawMessage(dsl.EachLike(string(b), len(data)))

	case json.RawMessage:
		return data

	case nil:
		return json.RawMessage("null")
	}

	panic(fmt.Sprintf("can not encode %T", data))
}
