package application

import (
	"context"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// LabelService for translatable labels
	LabelService struct {
		labelProvider              labelProvider
		defaultLocaleCode          string
		defaultFallbackLocaleCodes []string
		translationService         *infrastructure.TranslationService
	}

	labelProvider func() *domain.Label
)

// Inject dependencies
func (l *LabelService) Inject(labelProvider labelProvider, translationService *infrastructure.TranslationService, config *struct {
	DefaultLocaleCode string       `inject:"config:locale.locale"`
	FallbackLocalCode config.Slice `inject:"config:locale.fallbackLocales,optional"`
}) {
	l.translationService = translationService
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

// All labels for the API request
func (l *LabelService) AllLabels(ctx context.Context, r web.Request) []domain.Label {
	var labels []domain.Label
	tags := l.translationService.AllTranslationTags(l.defaultLocaleCode)

	for _, tag := range tags {
		label := l.NewLabel(tag)
		if label != nil {
			labels = append(labels, *label)
		}
	}

	return labels
}
