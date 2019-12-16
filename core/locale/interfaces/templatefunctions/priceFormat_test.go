package templatefunctions_test

import (
	"context"
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
)

type FakeTranslationService struct{}

func (s *FakeTranslationService) AllTranslationKeys(localeCode string) []string {
	return []string{
		"key1",
		"key2",
	}
}

func (s *FakeTranslationService) Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string {
	return defaultLabel
}

func (s *FakeTranslationService) TranslateLabel(label domain.Label) string {
	return label.GetDefaultLabel()
}

func FakeLabelProvider() *domain.Label {
	label := &domain.Label{}
	label.Inject(new(FakeTranslationService))
	return label
}

func TestPriceFormatFunc_Func(t *testing.T) {
	labelService := &application.LabelService{}

	labelService.Inject(FakeLabelProvider, nil, nil)

	type fields struct {
		config       config.Map `inject:"config:core.locale.accounting"`
		labelService *application.LabelService
	}
	type args struct {
		value    float64
		currency string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "Euro",
			fields: fields{
				config: config.Map{
					"default": config.Map{
						"decimal":    ",",
						"thousand":   ".",
						"formatZero": "%s -,-",
						"format":     "%s %v",
						"formatLong": "%v %v",
					},
				},
				labelService: labelService,
			},
			args: args{
				value:    21500.99,
				currency: "€",
			},
			want: "€ 21.500,99",
		},
		{
			name: "Dollar",
			fields: fields{
				config: config.Map{
					"default": config.Map{
						"decimal":    ".",
						"thousand":   ",",
						"formatZero": "%s -,-",
						"format":     "%s %v",
						"formatLong": "%v %v",
					},
				},
				labelService: labelService,
			},
			args: args{
				value:    55,
				currency: "$",
			},
			want: "$ 55.00",
		}, {
			name: "Dollar non default with no space",
			fields: fields{
				config: config.Map{
					"default": config.Map{
						"decimal":    ".",
						"thousand":   ",",
						"formatZero": "%s -,-",
						"format":     "%s %v",
						"formatLong": "%v %v",
					},
					"€": config.Map{
						"decimal":    ".",
						"thousand":   ",",
						"formatZero": "%s -,-",
						"format":     "%s%v",
						"formatLong": "%v %v",
					},
				},
				labelService: labelService,
			},
			args: args{
				value:    55,
				currency: "€",
			},
			want: "€55.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nff := &templatefunctions.PriceFormatFunc{}
			priceService := application.PriceService{}
			priceService.Inject(tt.fields.labelService, &struct {
				Config config.Map `inject:"config:core.locale.accounting"`
			}{tt.fields.config})
			nff.Inject(&priceService)

			templateFunc := nff.Func(context.Background()).(func(value float64, currency string) string)

			if got := templateFunc(tt.args.value, tt.args.currency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberFormatFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}
