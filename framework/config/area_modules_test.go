package config_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"flamingo.me/dingo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/config"
)

// parentModule returns a new ModuleFunc on every Depends call, which conflicts
// with the ModuleFunc of an earlier call if both are initialized.
type parentModule struct {
	dependsCalls int
	initialized  []string
}

func (m *parentModule) Depends() []dingo.Module {
	m.dependsCalls++
	call := m.dependsCalls

	return []dingo.Module{
		dingo.ModuleFunc(func(injector *dingo.Injector) {
			m.initialized = append(m.initialized, "func")

			injector.Bind(new(int)).AnnotatedWith("call").ToInstance(call)
		}),
		childModule{parent: m},
	}
}

func (m *parentModule) Configure(*dingo.Injector) {
	m.initialized = append(m.initialized, "parent")
}

// childModule is not a pointer, as disabling such modules must work too.
type childModule struct{ parent *parentModule }

func (m childModule) Configure(*dingo.Injector) {
	m.parent.initialized = append(m.parent.initialized, "child")
}

func (childModule) CueConfig() string { return `child: name: string | *"default"` }

func (childModule) DefaultConfig() config.Map { return config.Map{"child.fallback": "default"} }

func (childModule) OverrideConfig(config.Map) config.Map {
	return config.Map{"child.override": "default"}
}

func (childModule) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{"child.alias": "child.name"}
}

// captureLog also catches dingo's output, as slog writes through the log package.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()

	output := new(bytes.Buffer)
	previous := log.Writer()

	log.SetOutput(output)
	t.Cleanup(func() { log.SetOutput(previous) })

	return output
}

func TestAreaResolvesModulesOnce(t *testing.T) {
	t.Parallel()

	parent := new(parentModule)
	area := config.NewArea("root", []dingo.Module{parent})

	require.NoError(t, config.Load(area, t.TempDir()))
	_, err := area.GetInitializedInjector()
	require.NoError(t, err)

	assert.Equal(t, 1, parent.dependsCalls)
	assert.Equal(t, []string{"func", "child", "parent"}, parent.initialized)
}

func TestTryModulesResolvesModulesOnce(t *testing.T) {
	t.Parallel()

	parent := new(parentModule)

	require.NoError(t, config.TryModules(nil, parent))
	assert.Equal(t, 1, parent.dependsCalls)
}

//nolint:paralleltest // captures the global log output
func TestAreaDisabledModules(t *testing.T) {
	output := captureLog(t)
	parent := new(parentModule)
	area := config.NewArea("root", []dingo.Module{parent})
	dir := t.TempDir()

	// unknown.Module matches nothing, so it is compared with every module
	require.NoError(t, config.Load(area, dir, config.AdditionalConfig([]string{
		"flamingo.modules.disabled: [flamingo.me/flamingo/v3/framework/config_test.childModule, unknown.Module]",
	})))
	_, err := area.GetInitializedInjector()
	require.NoError(t, err)
	require.NoError(t, config.Load(area, dir))

	assert.Equal(t, []string{"func", "parent"}, parent.initialized, "a disabled dependency stays out")
	assert.Equal(t, 1, strings.Count(output.String(), "is disabled and not initialized"), output.String())
	assert.Contains(t, output.String(), "config_test.childModule is disabled")

	for _, key := range []string{"child.name", "child.fallback", "child.override", "child.alias"} {
		value, _ := area.Config(key)
		assert.Equal(t, "default", value, "a disabled module keeps its configuration: %s", key)
	}
}

type cycleModule struct{}

func (*cycleModule) Depends() []dingo.Module { return []dingo.Module{new(otherCycleModule)} }

func (*cycleModule) Configure(*dingo.Injector) {}

type otherCycleModule struct{}

func (*otherCycleModule) Depends() []dingo.Module { return []dingo.Module{new(cycleModule)} }

func (*otherCycleModule) Configure(*dingo.Injector) {}

func TestAreaDependencyCycle(t *testing.T) {
	t.Parallel()

	root := config.NewArea("root", []dingo.Module{new(cycleModule)})
	assert.ErrorIs(t, config.Load(root, t.TempDir()), dingo.ErrModuleCycle)

	parent := config.NewArea("root", nil, config.NewArea("child", []dingo.Module{new(cycleModule)}))
	require.NoError(t, config.Load(parent, t.TempDir()))

	_, err := parent.Flat()
	assert.ErrorIs(t, err, dingo.ErrModuleCycle)
}

type traced struct{}

type tracedModule struct{}

func (*tracedModule) Inject(*traced) {}

func (*tracedModule) Configure(*dingo.Injector) {}

//nolint:paralleltest // captures the global log output
func TestAreaInjectionTracing(t *testing.T) {
	output := captureLog(t)
	modules := []dingo.Module{new(tracedModule)}

	grandchild := config.NewArea("grandchild", modules)
	child := config.NewArea("child", modules, grandchild)
	sibling := config.NewArea("sibling", modules)
	root := config.NewArea("root", modules, child, sibling)
	child.SetDingoTracing(false, true)

	// parents come first, so each step initializes only the modules of one area
	for _, area := range []*config.Area{root, child, grandchild, sibling} {
		output.Reset()

		_, err := area.GetInitializedInjector()
		require.NoError(t, err)

		logged := strings.Contains(output.String(), "INJECTING: flamingo.me/flamingo/v3/framework/config_test#traced ")
		assert.Equal(t, area == child || area == grandchild, logged, area.Name)
	}
}

// loop constructs itself again through a provider until circular tracing or the
// call limit stops it.
type loop struct{}

type loopProvider func() *loop

type loopCalls struct{ count int }

func (*loop) Inject(next loopProvider, calls *loopCalls) {
	calls.count++
	if calls.count > 100 {
		panic("call limit reached")
	}

	next()
}

func TestAreaCircularTracing(t *testing.T) {
	t.Parallel()

	area := config.NewArea("root", []dingo.Module{dingo.ModuleFunc(func(injector *dingo.Injector) {
		injector.Bind(new(loopCalls)).ToInstance(new(loopCalls))
		assert.PanicsWithValue(t, "detected circular dependency", func() { _, _ = injector.GetInstance(new(loop)) })
	})})
	area.SetDingoTracing(true, false)

	_, err := area.GetInitializedInjector()
	require.NoError(t, err)
}
