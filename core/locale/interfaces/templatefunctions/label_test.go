package templatefunctions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

func TestLabelFormat_Func(t *testing.T) {
	translationService := &fake.TranslationService{}

	labelService := &application.LabelService{}
	labelService.Inject(func() *domain.Label {
		label := &domain.Label{}
		label.Inject(translationService)
		return label
	}, translationService, nil, &struct {
		DefaultLocaleCode string       `inject:"config:core.locale.locale"`
		FallbackLocalCode config.Slice `inject:"config:core.locale.fallbackLocales,optional"`
	}{
		DefaultLocaleCode: "en",
		FallbackLocalCode: config.Slice{"de"},
	})

	tFuncProvider := &LabelFunc{}
	tFuncProvider.Inject(labelService, flamingo.NullLogger{})

	tFunc, ok := tFuncProvider.Func(context.Background()).(func(key string, params ...interface{}) *domain.Label)
	assert.True(t, ok)

	result := tFunc("key", "deprecated")

	expected := &domain.Label{}
	expected.SetKey("key")
	expected.SetDefaultLabel("key")
	expected.SetLocaleCode("en")
	expected.SetFallbackLocaleCodes([]string{"de"})
	expected.SetCount(1)
	expected.Inject(&fake.TranslationService{})

	assert.Equal(t, expected, result)
}
