AuthHttp1 :: core.auth.http & {broker: "http1", realm: "http1 realm", users: {"user1": "pw1"}}
AuthHttp2 :: core.auth.http & {broker: "http2", realm: "http2 realm", users: {"user2": "pw2"}}
AuthFake1 :: core.auth.fake & {broker: "fake1", userConfig: {jondoe: {password: "password"}}}
AuthFake1WithoutPasswords :: core.auth.fake & {broker: "fake1WithoutPasswords", validatePassword: false, userConfig: {jondoe: {}}}
AuthFake1WithDefaultTemplate :: core.auth.fake & {broker: "fake1WithDefaultTemplate", validatePassword: true, validateOtp: true, userConfig: {jondoe: {password: "password", otp: "otp"}}}
AuthFake1WithCustomTemplate :: core.auth.fake & {broker: "fake1WithCustomTemplate", validatePassword: true, validateOtp: true, usernameFieldId: "customUsernameField", passwordFieldId: "customPasswordField", otpFieldId: "customOtpField", userConfig: {jondoe: {password: "password", otp: "otp"}}, loginTemplate: """
<body>
  <h1>Custom Login Template!</h1>
  <form name="fake-idp-form" action="{{.FormURL}}" method="post">
	<div>{{.Message}}</div>
	<label for="{{.UsernameID}}">Username</label>
	<input type="text" name="{{.UsernameID}}" id="{{.UsernameID}}">
	<label for="{{.PasswordID}}">Password</label>
  <input type="password" name="{{.PasswordID}}" id="{{.PasswordID}}">
	<label for="{{.OtpID}}">2 Factor OTP</label>
  <input type="text" name="{{.OtpID}}" id="{{.OtpID}}">
	<button type="submit" id="submit">Fake Login</button>
  </form>
</body>
"""}

core: auth: web: broker: [
	core.auth.oidc & {broker: "oidc1", clientID: "client1", clientSecret: "client1", "endpoint": "http://127.0.0.1:3351/dex"},
	AuthHttp1,
	core.auth.oidc & {broker: "oidc2", clientID: "client2", clientSecret: "client2", "endpoint": "http://127.0.0.1:3352/dex"},
	AuthHttp2,
	core.auth.oidc & {broker: "kc1", clientID: "client1", clientSecret: "", "endpoint": "http://127.0.0.1:3353/auth/realms/Realm1", enableOfflineToken: false},
	core.auth.oidc & {broker: "kc2", clientID: "client2", clientSecret: "", "endpoint": "http://127.0.0.1:3354/auth/realms/Realm2", enableOfflineToken: false},
	customOidcBroker,
	StaticAuthBroker & {broker: "static1", users: ["user1", "user2"]},
	AuthFake1,
	AuthFake1WithoutPasswords,
	AuthFake1WithDefaultTemplate,
	AuthFake1WithCustomTemplate,
]

flamingo: session: cookie: secure: false
