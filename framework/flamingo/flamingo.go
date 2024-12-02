package flamingo

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
)

var appVersion = "develop"

const (
	vcsRevisionSettingKey = "vcs.revision"
)

type (
	AppInfo struct {
		AppVersion      string
		VCSRevision     string
		RuntimeVersion  string
		MainPackagePath string
		FlamingoVersion string
	}
)

// AppVersion returns the application version
// set this during build with `go build -ldflags "-X flamingo.me/flamingo/v3/framework/flamingo.appVersion=1.2.3"`.
func AppVersion() string {
	return appVersion
}

// GetAppInfo provides basic application information like runtime version, flamingo version etc.
func GetAppInfo() AppInfo {
	appInfo := AppInfo{
		AppVersion:     AppVersion(),
		RuntimeVersion: runtime.Version(),
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		appInfo.MainPackagePath = info.Main.Path
		for _, module := range info.Deps {
			if module.Path == "flamingo.me/flamingo/v3" {
				appInfo.FlamingoVersion = module.Version
			}
		}

		for _, setting := range info.Settings {
			if setting.Key == vcsRevisionSettingKey {
				appInfo.VCSRevision = setting.Value
			}
		}
	}

	return appInfo

}

// PrintAppInfo prints application info to the writer
func PrintAppInfo(writer io.Writer, appInfo AppInfo) {
	_, _ = fmt.Fprintf(writer, "%20s\t%s\n", "App version:", appInfo.AppVersion)
	_, _ = fmt.Fprintf(writer, "%20s\t%s\n", "Go runtime version:", appInfo.RuntimeVersion)
	_, _ = fmt.Fprintf(writer, "%20s\t%s\n", "VCS revision:", appInfo.VCSRevision)
	_, _ = fmt.Fprintf(writer, "%20s\t%s\n", "Path:", appInfo.MainPackagePath)
	_, _ = fmt.Fprintf(writer, "%20s\t%s\n", "Flamingo version:", appInfo.FlamingoVersion)
}
