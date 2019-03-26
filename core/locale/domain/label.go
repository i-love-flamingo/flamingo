package domain

import "encoding/json"

type (
	// Label instance
	Label struct {
		key                  string
		defaultLabel         string
		localeCode           string
		fallbacklocaleCodes  []string
		count                int
		translationArguments map[string]interface{}
		translationService   TranslationService
	}
	// TranslationService defines the translation service
	TranslationService interface {
		Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string
		TranslateLabel(label Label) string
	}
)

//GetTranslationArguments for label
func (l *Label) GetTranslationArguments() map[string]interface{} {
	return l.translationArguments
}

//GetCount for label
func (l *Label) GetCount() int {
	return l.count
}

//GetKey for label
func (l *Label) GetKey() string {
	return l.key
}

//GetDefaultLabel for label
func (l *Label) GetDefaultLabel() string {
	return l.defaultLabel
}

//GetLocaleCode - for label
func (l *Label) GetLocaleCode() string {
	return l.localeCode
}

//GetFallbacklocaleCodes for label
func (l *Label) GetFallbacklocaleCodes() []string {
	return l.fallbacklocaleCodes
}

// Inject translation service
func (l *Label) Inject(translationService TranslationService) {
	l.translationService = translationService
}

//String implements fmt.Stringer - pinning to the non pointer by intent
func (l Label) String() string {
	return l.translationService.TranslateLabel(l)
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

// SetFallbackLocales on a label
func (l *Label) SetFallbackLocales(fallbackLocaleCodes []string) *Label {
	l.fallbacklocaleCodes = fallbackLocaleCodes
	return l
}

// NoFallbackLocales on a label - removes any fallback locale codes
func (l *Label) NoFallbackLocales() *Label {
	l.fallbacklocaleCodes = nil
	return l
}

// AddFallbackLocale on a label
func (l *Label) AddFallbackLocale(localeCode string) *Label {
	l.fallbacklocaleCodes = append(l.fallbacklocaleCodes, localeCode)
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
