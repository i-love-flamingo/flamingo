package interfaces

import (
	"net/http"
)

type serviceMock struct{}

func (*serviceMock) BaseURL() string {
	return "base"
}

func (*serviceMock) DetermineBase(_ *http.Request) string {
	return "determinedBase"
}

func (*serviceMock) BaseDomain() string {
	return "baseDomain"
}
