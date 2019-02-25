package interfaces

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/baseurl/application"
)

func TestCanonicalDomainFunc_Func(t *testing.T) {
	service := new(application.Service).Inject(&struct {
		BaseURL string `inject:"config:baseurl.url"`
		Scheme  string `inject:"config:baseurl.scheme"`
	}{BaseURL: "domain.base", Scheme: "http://"})

	fnc := new(CanonicalDomainFunc).Inject(service).Func(context.Background()).(func() string)
	got := fnc()
	want := "domain.base"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
