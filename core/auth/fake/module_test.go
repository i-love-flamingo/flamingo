package fake

import (
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
)

func TestModule(t *testing.T) {
	t.Parallel()
	if err := config.TryModules(config.Map{"flamingo.debug.mode": true}, new(Module)); err != nil {
		t.Error(err)
	}
}

func TestModule_CueConfig(t *testing.T) {
	t.Parallel()
	m := new(Module)

	assert.NotEmpty(t, m.CueConfig(), "module returns a cue config string")

	cueBuildInstance := build.NewContext().NewInstance("test", nil)
	assert.NoError(t, cueBuildInstance.AddFile("test", m.CueConfig()), "cue config parsed without error")

	cueInstance, err := new(cue.Runtime).Build(cueBuildInstance)
	assert.NoError(t, err, "test cue instance build without error")

	cueInstance, err = cueInstance.Fill(
		config.Map{
			"core": config.Map{
				"auth": config.Map{
					"fake": config.Map{
						"loginTemplate": "testTemplateName",
						"userConfig": config.Map{
							"testUserA": config.Map{
								"password": "testUserAPassword",
							},
							"testUserB": config.Map{
								"password": "testUserBPassword",
								"otp":      "testUserBOtp",
							},
						},
						"validatePassword": true,
						"validateOtp":      true,
						"usernameFieldId":  "username",
						"passwordFieldId":  "password",
						"otpFieldId":       "otp",
					},
				},
			},
		},
	)

	assert.NoError(t, cueInstance.Value().Validate(), "cue config not loadable")
}
