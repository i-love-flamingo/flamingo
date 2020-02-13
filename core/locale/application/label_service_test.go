package application

import (
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLabelService_NewLabel(t *testing.T) {
	provider := func() *domain.Label {
		label := &domain.Label{}
		label.SetTranslationArguments(map[string]interface{}{
			"argKey": "argValue",
		})
		return label
	}


	service := &LabelService{}
	service.Inject(provider, &fake.TranslationService{}, flamingo.NullLogger{}, &struct {
		DefaultLocaleCode string       `inject:"config:core.locale.locale"`
		FallbackLocalCode config.Slice `inject:"config:core.locale.fallbackLocales,optional"`
	}{
		DefaultLocaleCode: "en",
		FallbackLocalCode: config.Slice{"de"},
	})

	result := service.NewLabel("key")

	expected := &domain.Label{}
	expected.SetTranslationArguments(map[string]interface{}{
		"argKey": "argValue",
	})
	expected.SetKey("key")
	expected.SetDefaultLabel("key")
	expected.SetLocaleCode("en")
	expected.SetCount(1)
	expected.SetFallbackLocaleCodes([]string{"de"})

	assert.Equal(t, expected, result)
}

func TestLabelService_AllLabels(t *testing.T) {
	provider := func() *domain.Label {
		label := &domain.Label{}
		label.SetTranslationArguments(map[string]interface{}{
			"argKey": "argValue",
		})
		return label
	}


	service := &LabelService{}
	service.Inject(provider, &fake.TranslationService{}, flamingo.NullLogger{}, &struct {
		DefaultLocaleCode string       `inject:"config:core.locale.locale"`
		FallbackLocalCode config.Slice `inject:"config:core.locale.fallbackLocales,optional"`
	}{
		DefaultLocaleCode: "en",
		FallbackLocalCode: config.Slice{"de"},
	})

	result := service.AllLabels()
	var expected []domain.Label

	keys := []string{"key1", "key2"}
	for _, k := range keys {
		item := domain.Label{}
		item.SetTranslationArguments(map[string]interface{}{
			"argKey": "argValue",
		})
		item.SetKey(k)
		item.SetDefaultLabel(k)
		item.SetLocaleCode("en")
		item.SetCount(1)
		item.SetFallbackLocaleCodes([]string{"de"})

		expected = append(expected, item)
	}



	assert.Equal(t, expected, result)
}