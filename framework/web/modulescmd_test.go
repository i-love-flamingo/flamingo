package web_test

import (
	"testing"

	"flamingo.me/dingo"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	testModule1 struct{}
	testModule2 struct{}
	testModule3 struct{}
)

func (m1 *testModule1) Configure(_ *dingo.Injector) {}
func (m2 *testModule2) Configure(_ *dingo.Injector) {}
func (m3 *testModule3) Configure(_ *dingo.Injector) {}

func TestModulesCmd_Print(t *testing.T) {
	t.Run("", func(t *testing.T) {
		modules := []dingo.Module{
			new(testModule1),
			new(testModule2),
			new(testModule3),
		}
		testArea := config.NewArea("testArea", modules)

		function := web.ModulesCmd(testArea).Run
		assert.NotNil(t, function)
		function(nil, nil)
	})
}
