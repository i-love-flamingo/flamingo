package application

import (
	"fmt"

	"go.aoe.com/flamingo/framework/web"
)

// DataCollector for profiling
type DataCollector struct {
	AutoManager *AuthManager `inject:""`
}

// Collect data
func (dc *DataCollector) Collect(ctx web.Context) string {
	token, err := dc.AutoManager.IDToken(ctx)
	if err != nil {
		return "AuthManager: " + err.Error()
	}
	return fmt.Sprintf("AuthManager: %s issued from %s (issued at %s valid until %s)", token.Subject, token.Issuer, token.IssuedAt, token.Expiry)
}
