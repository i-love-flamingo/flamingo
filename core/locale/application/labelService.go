package application

import "flamingo.me/flamingo/v3/core/locale/domain"

type (
	// LabelService for translatable labels
	LabelService struct {
		labelProvider     labelProvider
		defaultLocaleCode string
	}

	labelProvider func() *domain.Label
)

// Inject dependencies
func (l *LabelService) Inject(labelProvider labelProvider, config *struct {
	DefaultLocaleCode string `inject:"config:locale.locale"`
}) {
	l.labelProvider = labelProvider
	l.defaultLocaleCode = config.DefaultLocaleCode
}

// NewLabel factory
func (l *LabelService) NewLabel(key string) *domain.Label {
	label := l.labelProvider()
	return label.SetKey(key).SetDefaultLabel(key).SetLocale(l.defaultLocaleCode).SetCount(1)
}
