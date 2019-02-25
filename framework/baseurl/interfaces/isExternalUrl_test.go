package interfaces

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/baseurl/application"
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

	service := new(application.Service).Inject(&struct {
		BaseURL string `inject:"config:baseurl.url"`
		Scheme  string `inject:"config:baseurl.scheme"`
	}{BaseURL: "baseDomain", Scheme: "http://"})

	fnc := new(IsExternalURL).Inject(service).Func(context.Background()).(func(string) bool)

	for _, tt := range tests {
		got := fnc(tt.url)
		if got != tt.want {
			t.Errorf("%q is %v, but should be %v", tt.url, got, tt.want)
		}
	}
}
