package responder_test

import (
	"strings"
	"testing"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/catalog"
	"flamingo.me/flamingo/v3/core/gotemplate"
)

// bootResponder starts an app with the catalog routes and the go template engine. The lines go into config.yml.
func bootResponder(tb testing.TB, lines ...string) *harness.App {
	tb.Helper()

	return harness.Boot(tb, []dingo.Module{new(gotemplate.Module), new(catalog.Module)}, harness.Config{
		Files: map[string]string{"config.yml": configYAML(lines...)},
	})
}

// configYAML gives a config.yml with the test templates and debug mode off, followed by the lines.
func configYAML(lines ...string) string {
	all := append([]string{
		"core.gotemplate.engine.templates.basepath: testdata/templates",
		"flamingo.debug.mode: false",
	}, lines...)

	return strings.Join(all, "\n") + "\n"
}
