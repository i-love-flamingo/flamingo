package flamingo

import (
	"fmt"
	"runtime/debug"
)

var appVersion = ""

const (
	baseSemanticVersion   = "v0.0.0"
	vcsRevisionSettingKey = "vcs.revision"
)

// AppVersion returns the application version
// set this during build with `go build -ldflags "-X flamingo.me/flamingo/v3/framework/flamingo.appVersion=1.2.3"`.
func AppVersion() string {
	if appVersion != "" {
		return appVersion
	}

	// in case no version is set with ldflags check executable build info (git commit hash is embedded by Go by default)
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == vcsRevisionSettingKey {
				return fmt.Sprintf("%s-%s", baseSemanticVersion, setting.Value[:8])
			}
		}
	}

	return baseSemanticVersion
}
