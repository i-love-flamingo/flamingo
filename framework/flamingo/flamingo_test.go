package flamingo_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

func TestPrintAppInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		appInfo flamingo.AppInfo
	}

	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name: "",
			args: args{
				appInfo: flamingo.AppInfo{
					AppVersion:      "v1.2.3",
					VCSRevision:     "c9ce01204a18ff2f3e9ed999fbf7f3eb8e70b614",
					RuntimeVersion:  "go1.23.3",
					MainPackagePath: "go.aoe.com/whitelabel-airline/flamingo",
					FlamingoVersion: "v3.11.0",
				},
			},
			wantWriter: "        App version:\tv1.2.3\n Go runtime version:\tgo1.23.3\n       VCS revision:\tc9ce01204a18ff2f3e9ed999fbf7f3eb8e70b614\n               Path:\tgo.aoe.com/whitelabel-airline/flamingo\n   Flamingo version:\tv3.11.0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			writer := &bytes.Buffer{}
			flamingo.PrintAppInfo(writer, tt.args.appInfo)
			assert.Equalf(t, tt.wantWriter, writer.String(), "PrintAppInfo(%v, %v)", writer, tt.args.appInfo)
		})
	}
}
