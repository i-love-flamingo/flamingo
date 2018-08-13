package interfaces

import (
	"testing"
)

func TestCanonicalDomainFunc_Func(t *testing.T) {
	fnc := new(CanonicalDomainFunc).Inject(new(applicationServiceMock)).Func().(func() string)
	got := fnc()
	want := new(applicationServiceMock).GetBaseDomain()

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
