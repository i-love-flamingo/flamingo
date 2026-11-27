//go:build !race

package perf_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

const (
	configLoadGrowthCeiling = 30506
	configLoadRuns          = 3
	configEnvCount          = 500
)

// The ceiling is for the growth from 0 to 500 environment variables: in the total count, the fixed cost of a
// load would hide a few percent more work per variable.
func TestPerf_AllocsConfigLoad500Env(t *testing.T) { //nolint:paralleltest // allocation counts need a test that is not parallel
	dir := t.TempDir() // before Clearenv, because TMPDIR is then gone

	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.yml"), []byte("core.zap.loglevel: Error\n"), 0o600))

	saved := os.Environ()

	os.Clearenv()

	t.Cleanup(func() {
		os.Clearenv()

		for _, entry := range saved {
			name, value, _ := strings.Cut(entry, "=")

			if err := os.Setenv(name, value); err != nil { //nolint:usetesting // the cleanup restores the ambient environment, t.Setenv cannot
				t.Errorf("restore %s: %v", name, err)
			}
		}
	})

	load := func() *config.Area {
		area := config.NewArea("root", []dingo.Module{new(framework.InitModule), new(zap.Module), new(flamingo.SessionModule)})

		if err := config.Load(area, dir); err != nil {
			t.Error(err)
		}

		return area
	}

	base := minAllocs(t, configLoadRuns, func() { load() })

	t.Logf("%v allocations per run with 0 environment variables", base)

	for i := range configEnvCount {
		require.NoError(t, os.Setenv(fmt.Sprintf("CONTRACT_PERF_VAR_%03d", i), "value")) //nolint:usetesting // the test must remove the ambient environment, t.Setenv cannot
	}

	withEnv := minAllocs(t, configLoadRuns, func() { load() })

	value, found := load().Configuration.Get("flamingo.os.env.CONTRACT_PERF_VAR_000")
	require.True(t, found, "the variables must reach the config")
	require.Equal(t, "value", value)

	t.Logf("%v allocations per run with %d environment variables", withEnv, configEnvCount)

	assertAllocs(t, withEnv-base, configLoadGrowthCeiling)
}
