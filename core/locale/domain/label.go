package domain

import "encoding/json"

// Label instance
type Label struct {
	key                  string
	defaultLabel         string
	localeCode           string
	fallbackLocaleCodes  []string
	count                int
	translationArguments map[string]interface{}
	translationService   TranslationService
}

// Inject translation service
func (l *Label) Inject(translationService TranslationService) {
	l.translationService = translationService
}

// GetTranslationArguments for label
func (l *Label) GetTranslationArguments() map[string]interface{} {
	return l.translationArguments
}

// GetCount for label
func (l *Label) GetCount() int {
	return l.count
}

// GetKey for label
func (l *Label) GetKey() string {
	return l.key
}

// GetDefaultLabel for label
func (l *Label) GetDefaultLabel() string {
	return l.defaultLabel
}

// GetLocaleCode - for label
func (l *Label) GetLocaleCode() string {
	return l.localeCode
}

// GetFallbacklocaleCodes for label
// Deprecated: use GetFallbackLocaleCodes instead
func (l *Label) GetFallbacklocaleCodes() []string {
	return l.GetFallbackLocaleCodes()
}

// GetFallbackLocaleCodes for label
func (l *Label) GetFallbackLocaleCodes() []string {
	return l.fallbackLocaleCodes
}

// String implements fmt.Stringer - pinning to the non pointer by intent
func (l Label) String() string {
	return l.translationService.TranslateLabel(l)
}

// MarshalJSON implements fmt.Stringer - pinning to the non pointer by intent
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
// Deprecated: use SetLocaleCode instead
func (l *Label) SetLocale(localeCode string) *Label {
	return l.SetLocaleCode(localeCode)
}

// SetLocaleCode on a label
func (l *Label) SetLocaleCode(localeCode string) *Label {
	l.localeCode = localeCode
	return l
}

// SetFallbackLocales on a label
// Deprecated: use SetFallbackLocaleCodes instead
func (l *Label) SetFallbackLocales(fallbackLocaleCodes []string) *Label {
	return l.SetFallbackLocaleCodes(fallbackLocaleCodes)
}

// SetFallbackLocaleCodes on a label
func (l *Label) SetFallbackLocaleCodes(fallbackLocaleCodes []string) *Label {
	l.fallbackLocaleCodes = fallbackLocaleCodes
	return l
}

// NoFallbackLocales on a label - removes any fallback locale codes
// Deprecated: use NoFallbackLocaleCodes instead
func (l *Label) NoFallbackLocales() *Label {
	return l.NoFallbackLocaleCodes()
}

// NoFallbackLocaleCodes on a label - removes any fallback locale codes
func (l *Label) NoFallbackLocaleCodes() *Label {
	l.fallbackLocaleCodes = nil
	return l
}

// AddFallbackLocale on a label
// Deprecated: use AddFallbackLocaleCode instead
func (l *Label) AddFallbackLocale(localeCode string) *Label {
	return l.AddFallbackLocaleCode(localeCode)
}

// AddFallbackLocaleCode on a label
func (l *Label) AddFallbackLocaleCode(localeCode string) *Label {
	l.fallbackLocaleCodes = append(l.fallbackLocaleCodes, localeCode)
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
