package interfaces

import (
	"context"
	"testing"
)

func TestCanonicalUrlFunc_Func(t *testing.T) {
	fnc := new(CanonicalURLFunc).Inject(new(applicationServiceMock)).Func(context.Background()).(func() string)
	got := fnc()
	want := new(applicationServiceMock).GetCanonicalURLForCurrentRequest(nil)

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
