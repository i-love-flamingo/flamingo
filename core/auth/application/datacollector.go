package application

import (
	"flamingo/framework/web"
	"fmt"
)

// DataCollector for profiling
type DataCollector struct {
	AutoManager *AuthManager `inject:""`
}

// Collect data
func (dc *DataCollector) Collect(ctx web.Context) string {
	token, err := dc.AutoManager.IdToken(ctx)
	if err != nil {
		return "AuthManager: " + err.Error()
	}
	return fmt.Sprintf("AuthManager: %s issued from %s (issued at %s valid until %s)", token.Subject, token.Issuer, token.IssuedAt, token.Expiry)
}
