package config_test

import (
	"testing"

	"flamingo.me/dingo"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/framework/config"
)

type (
	testModule1 struct{}
	testModule2 struct{}
	testModule3 struct{}
	testModule4 struct{}
	testModule5 struct{}
)

func (m1 *testModule1) Configure(_ *dingo.Injector) {}
func (m2 *testModule2) Configure(_ *dingo.Injector) {}
func (m3 *testModule3) Configure(_ *dingo.Injector) {}
func (m3 *testModule4) Configure(_ *dingo.Injector) {}
func (m3 *testModule5) Configure(_ *dingo.Injector) {}

func TestModulesCmd_Print(t *testing.T) {
	t.Run("Visual test: print modules without duplicates", func(t *testing.T) {
		modules := []dingo.Module{
			new(testModule1),
			new(testModule2),
			new(testModule3),
		}

		childModules := []dingo.Module{
			new(testModule1),
			new(testModule2),
			new(testModule3),
			new(testModule4),
		}

		subChildModules := []dingo.Module{
			new(testModule1),
			new(testModule2),
			new(testModule3),
			new(testModule4),
			new(testModule5),
		}

		testArea := config.NewArea("testArea", modules)
		childArea := config.NewArea("childTestArea", childModules)
		subChildArea := config.NewArea("subChildTestArea", subChildModules)
		childArea.Childs = []*config.Area{
			subChildArea,
		}
		testArea.Childs = []*config.Area{
			childArea,
		}

		function := config.ModulesCmd(testArea).Run
		assert.NotNil(t, function)
		function(nil, nil)
	})

	t.Run("Visual test: print modules with duplicates", func(t *testing.T) {
		modules := []dingo.Module{
			new(testModule1),
			new(testModule2),
			new(testModule3),
		}

		childModules := []dingo.Module{
			new(testModule1),
			new(testModule2),
			new(testModule3),
			new(testModule4),
		}

		subChildModules := []dingo.Module{
			new(testModule1),
			new(testModule2),
			new(testModule3),
			new(testModule4),
			new(testModule5),
		}

		testArea := config.NewArea("testArea", modules)
		childArea := config.NewArea("childTestArea", childModules)
		subChildArea := config.NewArea("subChildTestArea", subChildModules)
		childArea.Childs = []*config.Area{
			subChildArea,
		}
		testArea.Childs = []*config.Area{
			childArea,
		}

		function := config.ModulesCmd(testArea).Run
		assert.NotNil(t, function)
		function(nil, []string{"all"})
	})
}
