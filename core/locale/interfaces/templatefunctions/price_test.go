package templatefunctions_test

import (
	"context"
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
)

func fakeLabelProvider() *domain.Label {
	label := &domain.Label{}
	label.Inject(new(fake.TranslationService))
	return label
}

func TestPriceFormatFunc_Func(t *testing.T) {
	labelService := &application.LabelService{}

	labelService.Inject(fakeLabelProvider, nil, nil, nil)

	type fields struct {
		config       config.Map
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
			priceService.Inject(tt.fields.labelService, nil, &struct {
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

func TestPriceFormatLongFunc_Func(t *testing.T) {
	labelService := &application.LabelService{}

	labelService.Inject(fakeLabelProvider, nil, nil, nil)

	type fields struct {
		config       config.Map
		labelService *application.LabelService
	}
	type args struct {
		value         float64
		currency      string
		currencyLabel string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "$ USD",
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
				value:         21500.99,
				currency:      "$",
				currencyLabel: "USD",
			},
			want: "$ 21,500.99 USD",
		},
		{
			name: "$ USD no space",
			fields: fields{
				config: config.Map{
					"default": config.Map{
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
				value:         21500.99,
				currency:      "$",
				currencyLabel: "USD",
			},
			want: "$21,500.99 USD",
		},
		{
			// this testcase has switched the . and , for price format
			name: "$ NZD no space",
			fields: fields{
				config: config.Map{
					"default": config.Map{
						"decimal":    ".",
						"thousand":   ",",
						"formatZero": "%s -,-",
						"format":     "%s%v",
						"formatLong": "%v %v",
					},
					"$": config.Map{
						"decimal":    ",",
						"thousand":   ".",
						"formatZero": "%s -,-",
						"format":     "%s%v",
						"formatLong": "%v %v",
					},
				},
				labelService: labelService,
			},
			args: args{
				value:         21500.99,
				currency:      "$",
				currencyLabel: "NZD",
			},
			want: "$21.500,99 NZD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			priceService := application.PriceService{}
			priceService.Inject(tt.fields.labelService, nil, &struct {
				Config config.Map `inject:"config:core.locale.accounting"`
			}{tt.fields.config})

			priceFormatFunc := &templatefunctions.PriceFormatFunc{}
			priceFormatFunc.Inject(&priceService)

			priceFormatLongFunc := &templatefunctions.PriceFormatLongFunc{}
			priceFormatLongFunc.Inject(tt.fields.labelService, priceFormatFunc, &priceService, &struct {
				Config config.Map `inject:"config:core.locale.accounting"`
			}{tt.fields.config})

			templateFunc := priceFormatLongFunc.Func(context.Background()).(func(value float64, currency string, currencyLabel string) string)

			if got := templateFunc(tt.args.value, tt.args.currency, tt.args.currencyLabel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberFormatFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}
