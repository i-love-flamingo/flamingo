package fake

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule(t *testing.T) {
	t.Parallel()
	c := config.Map{
		"flamingo.debug.mode": true,
		"core.auth.fake": config.Map{
			"broker":        "fakeBroker",
			"loginTemplate": "testTemplateName",
			"userConfig": config.Map{
				"testUserA": config.Map{
					"password": "testUserAPassword",
				},
				"testUserB": config.Map{
					"password": "testUserBPassword",
				},
			},
			"validatePassword": true,
			"usernameFieldId":  "username",
			"passwordFieldId":  "password",
		},
	}
	if err := config.TryModules(c, new(Module)); err != nil {
		t.Error(err)
	}
}
