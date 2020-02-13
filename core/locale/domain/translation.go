package domain

// TranslationService defines the translation service
type TranslationService interface {
	Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string
	TranslateLabel(label Label) string
	AllTranslationKeys(localeCode string) []string
}
