package flamingo

var appVersion = "develop"

// AppVersion returns the application version
// set this during build with `go build -ldflags "-X flamingo.me/flamingo/v3/framework/flamingo.appVersion=1.2.3"`.
func AppVersion() string {
	return appVersion
}
