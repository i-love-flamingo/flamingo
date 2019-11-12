package controllers

import (
	"context"
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	TranslationController struct {
		responder    *web.Responder
		labelService *application.LabelService
	}

	TranslationJson struct {
		Key         string `json:"key"`
		Translation string `json:"translation"`
	}
)

func (c *TranslationController) Inject(
	responder *web.Responder,
	labelService *application.LabelService,
) {
	c.responder = responder
	c.labelService = labelService
}

func (c *TranslationController) GetAllTranslations(ctx context.Context, r *web.Request) web.Result {
	translations := []TranslationJson{}
	l := c.labelService.AllLabels(ctx, *r)

	for _, la := range l {
		translations = append(translations, TranslationJson{
			Key:         la.GetKey(),
			Translation: la.String(),
		})
	}

	return c.responder.Data(translations)
}
