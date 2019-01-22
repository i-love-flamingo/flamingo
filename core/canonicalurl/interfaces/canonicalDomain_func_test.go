package interfaces

import (
	"context"
	"testing"
)

func TestCanonicalDomainFunc_Func(t *testing.T) {
	fnc := new(CanonicalDomainFunc).Inject(new(applicationServiceMock)).Func(context.TODO()).(func() string)
	got := fnc()
	want := new(applicationServiceMock).GetBaseDomain()

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
