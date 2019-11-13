package controllers

import (
	"context"
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// TranslationController to be used to return translations for all labels as array
	TranslationController struct {
		responder    *web.Responder
		labelService *application.LabelService
	}

	// TranslationJSON helper struct to map the result
	TranslationJSON struct {
		Key         string `json:"key"`
		Translation string `json:"translation"`
	}
)

// Inject dependencies
func (c *TranslationController) Inject(
	responder *web.Responder,
	labelService *application.LabelService,
) {
	c.responder = responder
	c.labelService = labelService
}

// GetAllTranslations controller for TranslationController
func (c *TranslationController) GetAllTranslations(ctx context.Context, r *web.Request) web.Result {
	translations := []TranslationJSON{}
	l := c.labelService.AllLabels()

	for _, la := range l {
		translations = append(translations, TranslationJSON{
			Key:         la.GetKey(),
			Translation: la.String(),
		})
	}

	return c.responder.Data(translations)
}
