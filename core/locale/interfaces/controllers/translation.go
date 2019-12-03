package controllers

import (
	"context"
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/web"
	"sync"
)

type (
	// TranslationController to be used to return translations for all labels as array
	TranslationController struct {
		responder    *web.Responder
		labelService *application.LabelService
		labelCache   []TranslationJSON
		lock         sync.RWMutex
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
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.labelCache == nil {
		c.lock.RUnlock()
		c.lock.Lock()

		l := c.labelService.AllLabels()
		translations := make([]TranslationJSON, 0, len(l))

		for _, la := range l {
			translations = append(translations, TranslationJSON{
				Key:         la.GetKey(),
				Translation: la.String(),
			})
		}
		c.labelCache = translations
		c.lock.Unlock()
		c.lock.RLock()
	}

	return c.responder.Data(c.labelCache)
}
