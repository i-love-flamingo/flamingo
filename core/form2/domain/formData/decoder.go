package formdata

import (
	"context"
	"net/url"
	"reflect"
	"strings"

	"github.com/leebenson/conform"

	"github.com/go-playground/form"

	"flamingo.me/flamingo/v3/core/form2/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// DefaultFormDataDecoderImpl represents implementation of default domain.FormDataDecoder.
	DefaultFormDataDecoderImpl struct{}
)

var _ domain.DefaultFormDataDecoder = &DefaultFormDataDecoderImpl{}

// Decode performs default form data decoding, depending if passed form data is instance of map[string]string or any other interface.
func (p *DefaultFormDataDecoderImpl) Decode(_ context.Context, _ *web.Request, values url.Values, formData interface{}) (interface{}, error) {
	if _, ok := formData.(map[string]string); ok {
		return p.decodeStringMap(values), nil
	}

	return p.decodeUnknownInterface(values, formData)
}

// decodeStringMap performs form data decoding by storing all POST values into simple instance of map[string]string.
func (p *DefaultFormDataDecoderImpl) decodeStringMap(values url.Values) map[string]string {
	stringMap := make(map[string]string, len(values))

	if values == nil {
		return stringMap
	}

	for k, v := range values {
		stringMap[k] = strings.Join(v, " ")
	}

	return stringMap
}

// decodeUnknownInterface performs form data decoding by using decoder from go-playground form package.
// It also performs string values' optimization byt using conform package.
func (p *DefaultFormDataDecoderImpl) decodeUnknownInterface(values url.Values, formData interface{}) (interface{}, error) {
	typeOf := reflect.TypeOf(formData)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	zeroFormData := reflect.New(typeOf).Interface()

	if values == nil {
		values = url.Values{}
	}

	decoder := form.NewDecoder()
	err := decoder.Decode(&zeroFormData, values)
	if err != nil {
		return nil, err
	}

	err = conform.Strings(zeroFormData)
	if err != nil {
		return nil, err
	}

	finalFormData := reflect.ValueOf(zeroFormData)
	if finalFormData.Kind() == reflect.Ptr {
		return finalFormData.Elem().Interface(), nil
	}

	return zeroFormData, nil
}
