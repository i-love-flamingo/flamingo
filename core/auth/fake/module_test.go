package fake

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule(t *testing.T) {
	t.Parallel()
	c := config.Map{
		"flamingo.debug.mode": true,
		"fake": config.Map{
			"loginTemplate": "testTemplateName",
			"userConfig": config.Map{
				"testUserA": config.Map{
					"password": "testUserAPassword",
				},
				"testUserB": config.Map{
					"password": "testUserBPassword",
					"otp":      "testUserBotp",
				},
			},
			"validatePassword": true,
			"validateOtp":      true,
			"usernameFieldId":  "username",
			"passwordFieldId":  "password",
			"otpFieldId":       "otp",
		},
	}
	if err := config.TryModules(c, new(Module)); err != nil {
		t.Error(err)
	}
}
