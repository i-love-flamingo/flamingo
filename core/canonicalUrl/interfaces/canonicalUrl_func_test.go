package interfaces

import (
	"context"
	"testing"
)

func TestCanonicalUrlFunc_Func(t *testing.T) {
	fnc := new(CanonicalUrlFunc).Inject(new(applicationServiceMock)).Func(context.Background()).(func() string)
	got := fnc()
	want := new(applicationServiceMock).GetCanonicalUrlForCurrentRequest(nil)

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
