package templatefunctions_test

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

func TestNumberFormatFunc_Func(t *testing.T) {
	type fields struct {
		precision float64
		decimal   string
		thousand  string
		logger    flamingo.Logger
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "int",
			fields: fields{
				precision: 1,
				decimal:   ".",
				thousand:  ",",
				logger:    flamingo.NullLogger{},
			},
			args: args{
				value: 55,
			},
			want: "55.0",
		},
		{
			name: "float64",
			fields: fields{
				precision: 3,
				decimal:   ",",
				thousand:  "",
				logger:    flamingo.NullLogger{},
			},
			args: args{
				value: float64(55),
			},
			want: "55,000",
		},
		{
			name: "BigFloat",
			fields: fields{
				precision: 2,
				decimal:   ".",
				thousand:  ",",
				logger:    flamingo.NullLogger{},
			},
			args: args{
				value: big.NewFloat(5500.0),
			},
			want: "5,500.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nff := &templatefunctions.NumberFormatFunc{}
			nff.Inject(tt.fields.logger, &struct {
				Precision float64 `inject:"config:core.locale.numbers.precision"`
				Decimal   string  `inject:"config:core.locale.numbers.decimal"`
				Thousand  string  `inject:"config:core.locale.numbers.thousand"`
			}{
				Precision: tt.fields.precision,
				Decimal:   tt.fields.decimal,
				Thousand:  tt.fields.thousand,
			})

			templateFunc := nff.Func(context.Background()).(func(value interface{}, params ...int) string)

			if got := templateFunc(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberFormatFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}
