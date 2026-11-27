package bootstrap_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/ops"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

// bootOpsAreas starts an app with the ops module in root and the areas "en" and "de".
func bootOpsAreas(t *testing.T) *harness.App {
	t.Helper()

	return harness.Boot(t,
		[]dingo.Module{new(ops.Module), new(prefixrouter.Module), harness.PassThroughWrapper()},
		harness.Config{
			Files: areaFiles(),
			Areas: []*config.Area{
				config.NewArea("en", nil),
				config.NewArea("de", nil),
			},
		},
	)
}

// instanceID returns the instance ID that an ops route gives.
func instanceID(t *testing.T, app *harness.App, path string) string {
	t.Helper()

	body := app.Get(t, path).Body

	var data map[string]string

	require.NoError(t, json.Unmarshal([]byte(body), &data), body)
	require.NotEmpty(t, data["id"], body)

	return data["id"]
}

func TestScope_SingletonSharedAcrossAreas(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootOpsAreas(t)

	en := instanceID(t, app, "/en/ops/instance/singleton")
	de := instanceID(t, app, "/de/ops/instance/singleton")

	assert.Equal(t, en, de)
}

func TestScope_ChildSingletonOnePerArea(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootOpsAreas(t)

	en := instanceID(t, app, "/en/ops/instance/child")
	de := instanceID(t, app, "/de/ops/instance/child")

	assert.NotEqual(t, en, de)
	assert.Equal(t, en, instanceID(t, app, "/en/ops/instance/child"))
}
