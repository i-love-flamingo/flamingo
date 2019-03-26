package application

import (
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/framework/config"
)

type (
	// LabelService for translatable labels
	LabelService struct {
		labelProvider              labelProvider
		defaultLocaleCode          string
		defaultFallbackLocaleCodes []string
	}

	labelProvider func() *domain.Label
)

// Inject dependencies
func (l *LabelService) Inject(labelProvider labelProvider, config *struct {
	DefaultLocaleCode string       `inject:"config:locale.locale"`
	FallbackLocalCode config.Slice `inject:"config:locale.fallbackLocales,optional"`
}) {
	l.labelProvider = labelProvider
	if config != nil {
		l.defaultLocaleCode = config.DefaultLocaleCode
		config.FallbackLocalCode.MapInto(&l.defaultFallbackLocaleCodes)
	}
}

// NewLabel factory
func (l *LabelService) NewLabel(key string) *domain.Label {
	label := l.labelProvider()
	return label.SetKey(key).SetDefaultLabel(key).SetLocale(l.defaultLocaleCode).SetCount(1).SetFallbackLocales(l.defaultFallbackLocaleCodes)
}
