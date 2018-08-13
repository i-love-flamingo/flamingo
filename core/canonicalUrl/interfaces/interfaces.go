package interfaces

import (
	"context"
)

type ApplicationService interface {
	GetBaseDomain() string
	GetCanonicalUrlForCurrentRequest(context.Context) string
}
