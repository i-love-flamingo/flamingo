package domain

import "encoding/json"

type (
	// Label instance
	Label struct {
		key                  string
		defaultLabel         string
		localeCode           string
		count                int
		translationArguments map[string]interface{}
		translationService   TranslationService
	}
	// TranslationService defines the translation service
	TranslationService interface {
		Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string
	}
)

// Inject translation service
func (l *Label) Inject(translationService TranslationService) {
	l.translationService = translationService
}

//String implements fmt.Stringer - pinning to the non pointer by intent
func (l Label) String() string {
	return l.translationService.Translate(l.key, l.defaultLabel, l.localeCode, l.count, l.translationArguments)
}

//MarshalJSON implements fmt.Stringer - pinning to the non pointer by intent
func (l Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.translationService.Translate(l.key, l.defaultLabel, l.localeCode, l.count, l.translationArguments))
}

// SetTranslationArguments sets the argument map
func (l *Label) SetTranslationArguments(translationArguments map[string]interface{}) *Label {
	l.translationArguments = translationArguments
	return l
}

// SetCount on a label
func (l *Label) SetCount(count int) *Label {
	l.count = count
	return l
}

// SetLocale on a label
func (l *Label) SetLocale(localeCode string) *Label {
	l.localeCode = localeCode
	return l
}

// SetDefaultLabel on a label
func (l *Label) SetDefaultLabel(defaultLabel string) *Label {
	l.defaultLabel = defaultLabel
	return l
}

// SetKey on a label
func (l *Label) SetKey(key string) *Label {
	l.key = key
	return l
}
