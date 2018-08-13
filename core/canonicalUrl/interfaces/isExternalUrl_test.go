package interfaces

import (
	"testing"
)

func TestIsExternalUrl_Func(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"http://example.com/a", true},
		{"http://baseDomain/", false},
		{"-invalid", true},
		{"", true},
		{"a/b", true},
	}

	fnc := new(IsExternalUrl).Inject(new(applicationServiceMock)).Func().(func(string) bool)

	for _, tt := range tests {
		got := fnc(tt.url)
		if got != tt.want {
			t.Errorf("%q is %v, but should be %v", tt.url, got, tt.want)
		}
	}
}
