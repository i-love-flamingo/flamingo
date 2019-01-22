package application_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/captcha/application"
	"github.com/go-test/deep"
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

func TestGenerator_NewCaptcha(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Default",
			args: args{
				length: 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := setUpGenerator()

			got := g.NewCaptcha(tt.args.length)
			assert.NotEmpty(t, got.Hash, "hash is empty")
			assert.NotEmpty(t, got.Solution, "solution is empty")
			assert.Lenf(t, got.Solution, tt.args.length, "solution length is %d, expected %d", len(got.Solution), tt.args.length)

			proof, err := g.NewCaptchaByHash(got.Hash)
			assert.Nil(t, err)
			if diff := deep.Equal(got, proof); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestGenerator_NewCaptchaByHashAndSolution(t *testing.T) {
	type args struct {
		solution string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				solution: "1234567890",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := setUpGenerator()

			cSolution := g.NewCaptchaBySolution(tt.args.solution)
			cHash, err := g.NewCaptchaByHash(cSolution.Hash)

			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.NewCaptchaByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(cSolution, cHash); diff != nil {
				t.Error(diff)
			}
		})
	}
}
