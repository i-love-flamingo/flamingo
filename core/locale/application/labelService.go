package application

import "flamingo.me/flamingo/v3/core/locale/domain"

type (
	LabelService struct {
		labelProvider LabelProvider
		defaultLocaleCode string
	}

	LabelProvider func() *domain.Label
)

func (l *LabelService) Inject(labelProvider LabelProvider, config *struct {
	DefaultLocaleCode string `inject:"config:locale.locale"`
}) {
	l.labelProvider = labelProvider
	l.defaultLocaleCode = config.DefaultLocaleCode
}

func (l *LabelService) NewLabel(key string) *domain.Label {
	label := l.labelProvider()
	return label.SetKey(key).SetDefaultLabel(key).SetLocale(l.defaultLocaleCode).SetCount(1)
}
