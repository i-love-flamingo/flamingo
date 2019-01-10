package application_test

import (
	"reflect"
	"strings"
	"testing"

	"flamingo.me/flamingo/core/captcha/application"
)

func TestPseudoStore_Get(t *testing.T) {
	g := &application.Generator{}
	g.Inject(
		&struct {
			EncryptionPassPhrase string `inject:"config:captcha.encryptionPassPhrase"`
		}{
			EncryptionPassPhrase: "test",
		},
	)

	c := g.NewCaptchaBySolution("198822")

	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "simple decode case",
			args: args{
				id: c.Hash,
			},
			want: []byte{1, 9, 8, 8, 2, 2},
		},
		{
			name: "invalid base64",
			args: args{
				id: strings.TrimSuffix(c.Hash, "="),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &application.PseudoStore{}
			s.Inject(g)
			if got := s.Get(tt.args.id, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PseudoStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
