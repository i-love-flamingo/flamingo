/*
Domain package of the "form" module.

The overall purpose (bounded context) of the form module is to provide models and functionality around processing web formulars.
So this "domain" package contains core models and types for handling forms - as well as common validation funcs.

Added date and regex validators:

FormData struct {
	...
	DateOfBirth string `form:"dateOfBirth" validate:"required,dateformat,minimumage=18,maximumage=150"`
	Password 	string `form:"password" validate:"required,password"`
	...
}

Additional date validators: dateformat (default 2006-01-02), minimumage (by specifying age value)
and maximumage (by specifying age value).

Regex validators can be specified as part of "customRegex" map. Each name defines validator tag, each value defines
actual regex.

To setup up specific config, use:
form:
	validator:
		dateFormat: 02.01.2006
		customRegex:
			password: ^[a-z]*$

To initiate validator use:

type (
	validatorProvider func() *validator.Validate
	....
)
*/
package domain
