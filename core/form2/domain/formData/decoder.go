package formData

import (
	"context"
	"github.com/leebenson/conform"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/form"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	DefaultFormDataDecoderImpl struct{}
)

var _ domain.DefaultFormDataDecoder = &DefaultFormDataDecoderImpl{}

func (p *DefaultFormDataDecoderImpl) Decode(_ context.Context, _ *web.Request, values url.Values, formData interface{}) (interface{}, error) {
	if _, ok := formData.(map[string]string); ok {
		return p.decodeStringMap(values), nil
	}

	return p.decodeUnknownInterface(values, formData)
}

func (p *DefaultFormDataDecoderImpl) decodeStringMap(values url.Values) map[string]string {
	stringMap := map[string]string{}

	if values == nil {
		return stringMap
	}

	for k, v := range values {
		stringMap[k] = strings.Join(v, " ")
	}

	return stringMap
}

func (p *DefaultFormDataDecoderImpl) decodeUnknownInterface(values url.Values, formData interface{}) (interface{}, error) {
	zeroFormData := reflect.Zero(reflect.ValueOf(formData).Type()).Interface()

	if values == nil {
		return zeroFormData, nil
	}

	decoder := form.NewDecoder()
	err := decoder.Decode(&zeroFormData, values)
	if err != nil {
		return nil, err
	}

	err = conform.Strings(&zeroFormData)
	if err != nil {
		return nil, err
	}

	return zeroFormData, nil
}
