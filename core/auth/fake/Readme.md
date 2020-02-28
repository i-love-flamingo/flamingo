# Fake Auth Module

## Description

The module allows to use a fake login service without "real" IDP providers with flexible configurable user credentials.

The module also provides a customizable login form template and also allows you to provide your own login form. The form input is validation is also configurable.

## Configuration

```cue
core: auth: {
	...
	web: broker: [
		core.auth.fake & {broker: "fake1", userConfig: {jondoe: {password: "password"}}},
		...
		core.auth.fake & {broker: "fake1WithoutPasswords", validatePassword: false, userConfig: {jondoe: {}}},
		...
		core.auth.fake & {broker: "fake1WithDefaultTemplate", validatePassword: true, userConfig: {jondoe: {password: "password"}}},
		...
		core.auth.fake & {broker: "fake1WithCustomTemplate", validatePassword: true, usernameFieldId: "customUsernameField", passwordFieldId: "customPasswordField", userConfig: {jondoe: {password: "password"}}, loginTemplate: """
<body>
  <h1>Custom Login Template!</h1>
  <form name="fake-idp-form" action="{{.FormURL}}" method="post">
	<div>{{.Message}}</div>
	<label for="{{.UsernameID}}">Username</label>
	<input type="text" name="{{.UsernameID}}" id="{{.UsernameID}}">
	<label for="{{.PasswordID}}">Password</label>
    <input type="password" name="{{.PasswordID}}" id="{{.PasswordID}}">
	<button type="submit" id="submit">Fake Login</button>
  </form>
</body>
"""},
    ...
	]
}
```

Aside from the fake configuration you MUST define a unique broker id.

### Custom Template HTML

The custom template html must use the same go template placeholders as used in the field configuration:

```html
<body>
  <h1>Login!</h1>
  <form name="fake-idp-form" action="{{.FormURL}}" method="post">
	<div>{{.Message}}</div>
	<label for="{{.UsernameID}}">Username</label>
	<input type="text" name="{{.UsernameID}}" id="{{.UsernameID}}">
	<label for="{{.PasswordID}}">Password</label>
    <input type="password" name="{{.PasswordID}}" id="{{.PasswordID}}">
	<button type="submit" id="submit">Fake Login</button>
  </form>
</body>
```

The placeholders `{{.Message}}` and `{{.FormURL}}` are required, the form submit type `POST` is mandantory.
