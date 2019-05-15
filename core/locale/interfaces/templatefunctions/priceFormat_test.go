package templatefunctions_test

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
)

type FakeTranslationService struct{}

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

	labelService.Inject(FakeLabelProvider, nil)

	type fields struct {
		config       config.Map `inject:"config:locale.accounting"`
		labelService *application.LabelService
	}
	type args struct {
		value    interface{}
		currency string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "float64",
			fields: fields{
				config: config.Map{
					"decimal":    ",",
					"thousand":   ".",
					"formatZero": "%s -,-",
					"format":     "%s %v",
					"formatLong": "%v %v",
				},
				labelService: labelService,
			},
			args: args{
				value:    float64(21500.99),
				currency: "€",
			},
			want: "€ 21.500,99",
		},
		{
			name: "int",
			fields: fields{
				config: config.Map{
					"decimal":    ".",
					"thousand":   ",",
					"formatZero": "%s -,-",
					"format":     "%s %v",
					"formatLong": "%v %v",
				},
				labelService: labelService,
			},
			args: args{
				value:    int(55),
				currency: "$",
			},
			want: "$ 55.00",
		},
		{
			name: "big.Float",
			fields: fields{
				config: config.Map{
					"decimal":    ".",
					"thousand":   ",",
					"formatZero": "%s -,-",
					"format":     "%s %v",
					"formatLong": "%v %v",
				},
				labelService: labelService,
			},
			args: args{
				value:    big.NewFloat(21500.99),
				currency: "¥",
			},
			want: "¥ 21,500.99",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nff := &templatefunctions.PriceFormatFunc{}
			nff.Inject(tt.fields.labelService, &struct {
				Config config.Map `inject:"config:locale.accounting"`
			}{tt.fields.config})

			templateFunc := nff.Func(context.Background()).(func(value interface{}, currency string) string)

			if got := templateFunc(tt.args.value, tt.args.currency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberFormatFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}
