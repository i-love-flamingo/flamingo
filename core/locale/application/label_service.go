package application

import (
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// LabelService for translatable labels
type LabelService struct {
	labelProvider              labelProvider
	defaultLocaleCode          string
	defaultFallbackLocaleCodes []string
	translationService         domain.TranslationService
}

type labelProvider func() *domain.Label

// Inject dependencies
func (l *LabelService) Inject(labelProvider labelProvider, translationService domain.TranslationService, logger flamingo.Logger, config *struct {
	DefaultLocaleCode string       `inject:"config:core.locale.locale"`
	FallbackLocalCode config.Slice `inject:"config:core.locale.fallbackLocales,optional"`
}) {
	l.translationService = translationService
	l.labelProvider = labelProvider

	if config == nil {
		return
	}

	l.defaultLocaleCode = config.DefaultLocaleCode
	err := config.FallbackLocalCode.MapInto(&l.defaultFallbackLocaleCodes)
	if err != nil {
		logger.WithField("category", "LabelService").Warn(err)
	}
}

// NewLabel factory
func (l *LabelService) NewLabel(key string) *domain.Label {
	label := l.labelProvider()
	return label.SetKey(key).SetDefaultLabel(key).SetLocaleCode(l.defaultLocaleCode).SetCount(1).SetFallbackLocaleCodes(l.defaultFallbackLocaleCodes)
}

// AllLabels return a array of all labels
func (l *LabelService) AllLabels() []domain.Label {
	var labels []domain.Label
	tags := l.translationService.AllTranslationKeys(l.defaultLocaleCode)

	for _, tag := range tags {
		label := l.NewLabel(tag)
		if label != nil {
			labels = append(labels, *label)
		}
	}

	return labels
}
