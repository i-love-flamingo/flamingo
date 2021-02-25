package fake

import (
	"flamingo.me/flamingo/v3/core/locale/domain"
)

// TranslationService is fake TranslationService used for testing purposes
type TranslationService struct{}

// AllTranslationKeys returns simple list of keys
func (s *TranslationService) AllTranslationKeys(localeCode string) []string {
	return []string{
		"key1",
		"key2",
	}
}

// Translate returns simple translation result
func (s *TranslationService) Translate(_ string, defaultLabel string, _ string, _ int, _ map[string]interface{}) string {
	return defaultLabel
}

// TranslateLabel returns simple label result
func (s *TranslationService) TranslateLabel(label domain.Label) string {
	return label.GetDefaultLabel()
}
