package interfaces

import (
	"context"
)

// ApplicationService interface for general canonicalURL usage
type ApplicationService interface {
	GetBaseDomain() string
	GetCanonicalURLForCurrentRequest(context.Context) string
}
