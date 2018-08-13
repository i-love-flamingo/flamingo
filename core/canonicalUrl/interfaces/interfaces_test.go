package interfaces

import "context"

type applicationServiceMock struct{}

func (*applicationServiceMock) GetBaseDomain() string {
	return "baseDomain"
}

func (*applicationServiceMock) GetCanonicalUrlForCurrentRequest(context.Context) string {
	return "canonicalUrl"
}
