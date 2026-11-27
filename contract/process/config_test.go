package process_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/contract/internal/harness"
)

func TestConfig_InvalidValueExitsWithKeyPath(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct{ flag, path string }{
		{`ops.greeting: 5`, "ops.greeting"},
		{`flamingo.session.backend: "nosuch"`, "flamingo.session.backend"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()

			p := harness.Process(t, harness.Options{
				Main: testappMain,
				Dir:  "../internal/testapp",
				Args: []string{"--flamingo-config", tc.flag, "serve", "--addr", "127.0.0.1:0"},
			})

			assert.Equal(t, 1, p.Wait(t), tc.flag)
			assert.True(t, containsLine(p.Lines(), tc.path), tc.flag)
		})
	}
}
