package interfaces

import (
	"context"
	"testing"
)

func TestCanonicalDomainFunc_Func(t *testing.T) {
	fnc := new(CanonicalDomainFunc).Inject(new(serviceMock)).Func(context.Background()).(func() string)
	got := fnc()
	want := new(serviceMock).BaseDomain()

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
