package templatefunctions_test

import (
	"context"
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestPriceFormatLongFunc_Func(t *testing.T) {
	labelService := &application.LabelService{}

	labelService.Inject(FakeLabelProvider, nil)

	type fields struct {
		config       config.Map `inject:"config:locale.accounting"`
		labelService *application.LabelService
	}
	type args struct {
		value         interface{}
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
			name: "float64",
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
				value:         float64(21500.99),
				currency:      "$",
				currencyLabel: "USD",
			},
			want: "$ 21,500.99 USD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priceFormatFunc := &templatefunctions.PriceFormatFunc{}
			priceFormatFunc.Inject(tt.fields.labelService, &struct {
				Config config.Map `inject:"config:locale.accounting"`
			}{tt.fields.config})

			priceFormatLongFunc := &templatefunctions.PriceFormatLongFunc{}
			priceFormatLongFunc.Inject(tt.fields.labelService, priceFormatFunc, &struct {
				Config config.Map `inject:"config:locale.accounting"`
			}{tt.fields.config})

			templateFunc := priceFormatLongFunc.Func(context.Background()).(func(value interface{}, currency string, currencyLabel string) string)

			if got := templateFunc(tt.args.value, tt.args.currency, tt.args.currencyLabel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberFormatFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}
