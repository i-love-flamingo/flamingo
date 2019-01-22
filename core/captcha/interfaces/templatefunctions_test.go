package interfaces_test

import (
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/core/captcha/application"
	"flamingo.me/flamingo/v3/core/captcha/domain"
	"flamingo.me/flamingo/v3/core/captcha/interfaces"
	"github.com/stretchr/testify/assert"
)

func setUpGenerator() *application.Generator {
	g := &application.Generator{}
	g.Inject(
		&struct {
			EncryptionPassPhrase string `inject:"config:captcha.encryptionPassPhrase"`
		}{
			EncryptionPassPhrase: "test",
		},
	)

	return g
}

func TestCaptchaFunc_Func(t *testing.T) {
	type fields struct {
		len int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Default",
			fields: fields{
				len: 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &interfaces.CaptchaFunc{}
			f.Inject(&struct {
				Len float64 `inject:"config:captcha.len"`
			}{Len: float64(tt.fields.len)}, setUpGenerator())
			tf := f.Func()

			got := tf.(func() *domain.Captcha)()
			assert.Len(t, got.Solution, tt.fields.len)
		})
	}
}

func TestCaptchaImgFunc_Func(t *testing.T) {
	tests := []struct {
		name    string
		options []bool
		captcha *domain.Captcha
		want    string
	}{
		{
			name: "Default",
			captcha: &domain.Captcha{
				Solution: "1234",
				Hash:     "h1234h",
			},
			want: "/captcha/h1234h.png",
		},
		{
			name:    "Download",
			options: []bool{true},
			captcha: &domain.Captcha{
				Solution: "1234",
				Hash:     "h1234h",
			},
			want: "/captcha/download/h1234h.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &interfaces.CaptchaImgFunc{}
			if got := f.Func().(func(*domain.Captcha, ...bool) string)(tt.captcha, tt.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CaptchaImgFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCaptchaSoundFunc_Func(t *testing.T) {
	tests := []struct {
		name    string
		options []bool
		captcha *domain.Captcha
		want    string
	}{
		{
			name: "Default",
			captcha: &domain.Captcha{
				Solution: "1234",
				Hash:     "h1234h",
			},
			want: "/captcha/h1234h.wav",
		},
		{
			name:    "Download",
			options: []bool{true},
			captcha: &domain.Captcha{
				Solution: "1234",
				Hash:     "h1234h",
			},
			want: "/captcha/download/h1234h.wav",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &interfaces.CaptchaSoundFunc{}
			if got := f.Func().(func(*domain.Captcha, ...bool) string)(tt.captcha, tt.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CaptchaImgFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}
