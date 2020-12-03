package controllers

import (
	"context"
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslationController_GetAllTranslations(t *testing.T) {
	translationService := &fake.TranslationService{}

	labelService := &application.LabelService{}
	labelService.Inject(func() *domain.Label {
		label := &domain.Label{}
		label.Inject(translationService)
		return label
	}, translationService, nil, &struct {
		DefaultLocaleCode string       `inject:"config:core.locale.locale"`
		FallbackLocalCode config.Slice `inject:"config:core.locale.fallbackLocales,optional"`
	}{DefaultLocaleCode: "en", FallbackLocalCode: config.Slice{"de"}})

	controller := &TranslationController{}
	controller.Inject(&web.Responder{}, labelService)

	result := controller.GetAllTranslations(context.Background(), web.CreateRequest(nil, nil))
	response, ok := result.(*web.DataResponse)
	assert.True(t, ok)
	assert.Equal(t, []TranslationJSON{
		{
			Key:         "key1",
			Translation: "key1",
		},
		{
			Key:         "key2",
			Translation: "key2",
		},
	}, response.Data)
}
