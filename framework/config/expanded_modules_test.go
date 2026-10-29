package config_test

import (
	"errors"
	"fmt"
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type capturedConfigModule struct {
	dependsCalls int
	configured   []int
	value        int
}

func (module *capturedConfigModule) Depends() []dingo.Module {
	module.dependsCalls++
	value := module.dependsCalls

	return []dingo.Module{
		dingo.ModuleFunc(func(injector *dingo.Injector) {
			module.configured = append(module.configured, value)
			injector.Bind(new(int)).ToInstance(value)
		}),
		&capturedConfigDependency{value: value},
	}
}

func (module *capturedConfigModule) Inject(value int) { module.value = value }

func (*capturedConfigModule) Configure(*dingo.Injector) {}

type capturedConfigDependency struct {
	value int
}

func (*capturedConfigDependency) Configure(*dingo.Injector) {}

func (module *capturedConfigDependency) CueConfig() string {
	return fmt.Sprintf(`snapshot: cue: *"generation %d" | string`, module.value)
}

func (module *capturedConfigDependency) DefaultConfig() config.Map {
	return config.Map{"snapshot.default": fmt.Sprintf("generation %d", module.value)}
}

func (*capturedConfigDependency) OverrideConfig(config.Map) config.Map {
	return config.Map{"snapshot.override": "overridden"}
}

func TestAreaReusesExpandedModulesAcrossLoads(t *testing.T) {
	t.Parallel()

	module := &capturedConfigModule{}
	area := config.NewArea("snapshot", []dingo.Module{module})
	basedir := t.TempDir()
	debugCalls := 0

	require.NoError(t, config.Load(area, basedir, config.CueDebug([]string{"snapshot", "cue"}, func(data []byte, err error) {
		require.NoError(t, err)
		assert.Contains(t, string(data), "generation 1")

		debugCalls++
	})))
	require.NoError(t, config.Load(area, basedir))
	assert.Equal(t, 1, debugCalls)
	assert.Equal(t, 1, module.dependsCalls, "each area must resolve dependencies once")
	require.Len(t, area.Modules, 3)
	assert.Contains(t, area.Modules, module)

	for key, want := range map[string]string{
		"snapshot.cue":      "generation 1",
		"snapshot.default":  "generation 1",
		"snapshot.override": "overridden",
	} {
		value, ok := area.Config(key)
		assert.True(t, ok)
		assert.Equal(t, want, value)
	}

	_, err := area.GetInitializedInjector()
	require.NoError(t, err)
	assert.Equal(t, 1, module.dependsCalls)
	assert.Equal(t, []int{1}, module.configured)
	assert.Equal(t, 1, module.value)
}

func TestTryModulesUsesExpandedClosureOnce(t *testing.T) {
	t.Parallel()

	module := &capturedConfigModule{}
	require.NoError(t, config.TryModules(nil, module))
	assert.Equal(t, 1, module.dependsCalls)
	assert.Equal(t, []int{1}, module.configured)
	assert.Equal(t, 1, module.value)
}

func TestTryModulesRecoversConfigurationPanics(t *testing.T) {
	t.Parallel()

	for _, value := range []any{errors.New("configuration failed"), "configuration failed"} {
		t.Run(fmt.Sprintf("%T", value), func(t *testing.T) {
			t.Parallel()

			var err error

			require.NotPanics(t, func() {
				err = config.TryModules(nil, dingo.ModuleFunc(func(*dingo.Injector) {
					panic(value)
				}))
			})
			assert.ErrorContains(t, err, "configuration failed")

			if panicError, ok := value.(error); ok {
				assert.ErrorIs(t, err, panicError)
			}
		})
	}
}

type disabledConfigModule struct {
	configured *int
}

func (module disabledConfigModule) Configure(*dingo.Injector) { *module.configured++ }

type disabledDependencyModule struct {
	dependency dingo.Module
	configured *[]string
}

func (module *disabledDependencyModule) Depends() []dingo.Module {
	return []dingo.Module{
		dingo.ModuleFunc(func(*dingo.Injector) {
			*module.configured = append(*module.configured, "closure")
		}),
		module.dependency,
	}
}

func (module *disabledDependencyModule) Configure(*dingo.Injector) {
	*module.configured = append(*module.configured, "parent")
}

func TestAreaKeepsNamedDependencyDisabledAfterClosure(t *testing.T) {
	t.Parallel()

	for _, pointer := range []bool{true, false} {
		t.Run(fmt.Sprintf("pointer=%t", pointer), func(t *testing.T) {
			t.Parallel()

			disabledCalls := 0

			var configured []string

			var dependency dingo.Module = disabledConfigModule{configured: &disabledCalls}
			if pointer {
				dependency = &disabledConfigModule{configured: &disabledCalls}
			}

			module := &disabledDependencyModule{dependency: dependency, configured: &configured}
			area := config.NewArea("disabled", []dingo.Module{module})
			require.NoError(t, config.Load(area, t.TempDir(), config.AdditionalConfig([]string{
				"flamingo.modules.disabled: [flamingo.me/flamingo/v3/framework/config_test.disabledConfigModule]",
			})))

			var err error

			require.NotPanics(t, func() { _, err = area.GetInitializedInjector() })
			require.NoError(t, err)
			assert.Zero(t, disabledCalls, "initialization must not restore a disabled dependency")
			assert.Equal(t, []string{"closure", "parent"}, configured)
		})
	}
}

type retainedConfigDependency struct {
	generation int
	configured int
}

func (module *retainedConfigDependency) Configure(*dingo.Injector) { module.configured++ }

func (module *retainedConfigDependency) CueConfig() string {
	return fmt.Sprintf(`retained: cue: *"generation %d" | string`, module.generation)
}

func (module *retainedConfigDependency) DefaultConfig() config.Map {
	return config.Map{
		"retained.default": fmt.Sprintf("generation %d", module.generation),
		"retained.old":     "legacy value",
	}
}

func (*retainedConfigDependency) OverrideConfig(current config.Map) config.Map {
	value, _ := current.Get("retained.default")
	return config.Map{"retained.override": fmt.Sprintf("override %s", value)}
}

func (*retainedConfigDependency) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{"retained.old": "retained.alias"}
}

type retainedConfigModule struct {
	dependsCalls int
	closureCalls int
	configured   int
	value        int
	dependency   *retainedConfigDependency
}

func (module *retainedConfigModule) Depends() []dingo.Module {
	module.dependsCalls++
	value := module.dependsCalls
	module.dependency = &retainedConfigDependency{generation: value}

	return []dingo.Module{
		dingo.ModuleFunc(func(injector *dingo.Injector) {
			module.closureCalls++

			injector.Bind(new(int)).ToInstance(value)
		}),
		module.dependency,
	}
}

func (module *retainedConfigModule) Inject(value int) { module.value = value }

func (module *retainedConfigModule) Configure(*dingo.Injector) { module.configured++ }

func TestAreaRetainsDisabledModuleConfigurationAcrossReloads(t *testing.T) {
	t.Parallel()

	module := &retainedConfigModule{}
	area := config.NewArea("retained", []dingo.Module{module})
	basedir := t.TempDir()
	want := map[string]string{
		"retained.cue":      "generation 1",
		"retained.default":  "generation 1",
		"retained.override": "override generation 1",
		"retained.old":      "legacy value",
		"retained.alias":    "legacy value",
	}
	assertConfiguration := func() {
		for key, expected := range want {
			value, ok := area.Config(key)
			assert.True(t, ok, key)
			assert.Equal(t, expected, value, key)
		}
	}

	require.NoError(t, config.Load(area, basedir, config.AdditionalConfig([]string{
		"flamingo.modules.disabled: [flamingo.me/flamingo/v3/framework/config_test.retainedConfigDependency]",
	})))
	assertConfiguration()

	// MergeFrom returns a shallow copy; initializing one area must preserve the other's list.
	copied := config.MergeFrom(*area, config.Area{})
	dependency := module.dependency
	injector, err := area.GetInitializedInjector()
	require.NoError(t, err)
	assertConfiguration()
	require.Len(t, area.Modules, 2)
	assert.NotContains(t, area.Modules, dependency)
	assert.Same(t, module, area.Modules[1])
	require.Len(t, copied.Modules, 3)
	assert.Same(t, dependency, copied.Modules[1])
	assert.Same(t, module, copied.Modules[2])

	require.NoError(t, config.Load(area, basedir))
	assertConfiguration()
	require.Len(t, area.Modules, 2, "reload must preserve the filtered diagnostic list")
	assert.NotContains(t, area.Modules, dependency)
	assert.Same(t, dependency, copied.Modules[1])
	assert.Equal(t, 1, module.dependsCalls)
	assert.Equal(t, 1, module.closureCalls)
	assert.Equal(t, 1, module.configured)
	assert.Equal(t, 1, module.value)
	assert.Zero(t, dependency.configured)

	for key, expected := range want {
		value, err := injector.GetAnnotatedInstance(new(string), "config:"+key)
		require.NoError(t, err)
		assert.Equal(t, expected, value, key)
	}
}

type retryConfigModule struct {
	cyclic       bool
	dependsCalls int
}

func (*retryConfigModule) Configure(*dingo.Injector) {}

func (module *retryConfigModule) Depends() []dingo.Module {
	module.dependsCalls++
	if module.cyclic {
		return []dingo.Module{module}
	}

	return nil
}

func TestAreaRetriesFailedExpansion(t *testing.T) {
	t.Parallel()

	module := &retryConfigModule{cyclic: true}
	area := config.NewArea("retry", []dingo.Module{module})
	basedir := t.TempDir()

	var err error

	require.NotPanics(t, func() { err = config.Load(area, basedir) })
	require.ErrorIs(t, err, dingo.ErrModuleCycle)
	assert.ErrorContains(t, err, "retry:")

	module.cyclic = false

	require.NoError(t, config.Load(area, basedir))
	require.NoError(t, config.Load(area, basedir))
	assert.Equal(t, 2, module.dependsCalls, "retry a failed resolution, then retain the successful list")
}

func TestTryModulesPreservesCycleError(t *testing.T) {
	t.Parallel()

	var err error

	require.NotPanics(t, func() { err = config.TryModules(nil, &retryConfigModule{cyclic: true}) })
	assert.ErrorIs(t, err, dingo.ErrModuleCycle)
}
