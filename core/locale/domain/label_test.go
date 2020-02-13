package domain

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestTranslationService struct{}

func (s *TestTranslationService) AllTranslationKeys(localeCode string) []string {

	return []string{
		fmt.Sprintf("key1_%s", localeCode),
		fmt.Sprintf("key2_%s", localeCode),
	}
}

func (s *TestTranslationService) Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string {
	return fmt.Sprintf("%s_%s_%s_%d_%v", key, defaultLabel, localeCode, count, translationArguments)
}

func (s *TestTranslationService) TranslateLabel(label Label) string {
	return fmt.Sprintf("%s_%s", label.GetKey(), label.GetDefaultLabel())
}

func TestLabel_AddFallbackLocaleCode(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.fallbackLocaleCodes)

	label.AddFallbackLocaleCode("de")
	assert.Equal(t, []string{"de"}, label.fallbackLocaleCodes)

	label.AddFallbackLocaleCode("en")
	assert.Equal(t, []string{"de", "en"}, label.fallbackLocaleCodes)
}

func TestLabel_GetCount(t *testing.T) {
	label := &Label{}
	assert.Equal(t, 0, label.GetCount())

	label = &Label{
		count: 10,
	}
	assert.Equal(t, 10, label.GetCount())
}

func TestLabel_GetDefaultLabel(t *testing.T) {
	label := &Label{}
	assert.Equal(t, "", label.GetDefaultLabel())

	label = &Label{
		defaultLabel: "en",
	}
	assert.Equal(t, "en", label.GetDefaultLabel())
}

func TestLabel_GetFallbackLocaleCodes(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.GetFallbackLocaleCodes())

	label = &Label{
		fallbackLocaleCodes: []string{"de", "en"},
	}
	assert.Equal(t, []string{"de", "en"}, label.GetFallbackLocaleCodes())
}

func TestLabel_GetKey(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.GetKey())

	label = &Label{
		key: "key",
	}
	assert.Equal(t, "key", label.GetKey())
}

func TestLabel_GetLocaleCode(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.GetLocaleCode())

	label = &Label{
		localeCode: "code",
	}
	assert.Equal(t, "code", label.GetLocaleCode())
}

func TestLabel_GetTranslationArguments(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.GetTranslationArguments())

	label = &Label{
		translationArguments: map[string]interface{}{
			"key": "value",
		},
	}
	assert.Equal(t, map[string]interface{}{
		"key": "value",
	}, label.GetTranslationArguments())
}

func TestLabel_Inject(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.translationService)

	label.Inject(&TestTranslationService{})
	assert.Equal(t, &TestTranslationService{}, label.translationService)
}

func TestLabel_MarshalJSON(t *testing.T) {
	service := &TestTranslationService{}
	label := &Label{
		key:          "key",
		defaultLabel: "defaultLabel",
		localeCode:   "localeCode",
		count:        1,
		translationArguments: map[string]interface{}{
			"key": "value",
		},
	}
	label.Inject(service)

	jResult, jErr := json.Marshal(service.Translate("key", "defaultLabel", "localeCode", 1, map[string]interface{}{
		"key": "value",
	}))
	assert.NoError(t, jErr)

	lResult, lErr := label.MarshalJSON()
	assert.NoError(t, lErr)

	assert.Equal(t, jResult, lResult)
}

func TestLabel_NoFallbackLocaleCodes(t *testing.T) {
	label := &Label{
		fallbackLocaleCodes: []string{"de", "en"},
	}
	assert.Equal(t, []string{"de", "en"}, label.fallbackLocaleCodes)

	label.NoFallbackLocaleCodes()
	assert.Empty(t, label.fallbackLocaleCodes)
}

func TestLabel_SetCount(t *testing.T) {
	label := &Label{}
	assert.Equal(t, 0, label.count)

	label.SetCount(10)
	assert.Equal(t, 10, label.count)
}

func TestLabel_SetDefaultLabel(t *testing.T) {
	label := &Label{}
	assert.Equal(t, "", label.defaultLabel)

	label.SetDefaultLabel("en")
	assert.Equal(t, "en", label.defaultLabel)
}

func TestLabel_SetFallbackLocaleCodes(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.fallbackLocaleCodes)

	label.SetFallbackLocaleCodes([]string{"de", "en"})
	assert.Equal(t, []string{"de", "en"}, label.fallbackLocaleCodes)
}

func TestLabel_SetKey(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.key)

	label.SetKey("key")
	assert.Equal(t, "key", label.key)
}

func TestLabel_SetLocaleCode(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.localeCode)

	label.SetLocaleCode("code")
	assert.Equal(t, "code", label.localeCode)
}

func TestLabel_SetTranslationArguments(t *testing.T) {
	label := &Label{}
	assert.Empty(t, label.translationArguments)

	label.SetTranslationArguments(map[string]interface{}{
		"key": "value",
	})
	assert.Equal(t, map[string]interface{}{
		"key": "value",
	}, label.translationArguments)
}

func TestLabel_String(t *testing.T) {
	service := &TestTranslationService{}
	label := Label{
		key:          "key",
		defaultLabel: "defaultLabel",
		localeCode:   "localeCode",
		count:        1,
		translationArguments: map[string]interface{}{
			"key": "value",
		},
	}
	label.Inject(service)

	assert.Equal(t, service.TranslateLabel(label), label.String())
}
