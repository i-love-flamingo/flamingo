package templatefunctions

import (
	"context"

	"fmt"

	"flamingo.me/flamingo/v3/core/csrfPreventionFilter/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/session"
)

type (
	CsrfInputFunc struct {
		service application.Service
		logger  flamingo.Logger
	}
)

func (f *CsrfInputFunc) Inject(s application.Service, l flamingo.Logger) {
	f.service = s
	f.logger = l
}

func (f *CsrfInputFunc) Func(ctx context.Context) interface{} {
	return func() interface{} {
		s, ok := session.FromContext(ctx)
		if !ok {
			f.logger.WithField("csrf", "templateFunc").Error("can't find session")
			return ""
		}

		return fmt.Sprintf(`<input type="hidden" name="%s" value="%s" />`, application.TokenName, f.service.Generate(s))
	}
}
