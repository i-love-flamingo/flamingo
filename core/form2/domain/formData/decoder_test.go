package formData

import (
	"github.com/stretchr/testify/suite"
	"net/url"
	"testing"
)

type (
	DefaultFormDataDecoderImplTestSuite struct {
		suite.Suite

		decoder *DefaultFormDataDecoderImpl
	}

	formDataDecoderTestData struct {
		Text   string    `form:"text" conform:"trim"`
		Number int       `form:"number"`
		Slice  []float64 `form:"slice"`
	}
)

func TestDefaultFormDataDecoderImplTestSuite(t *testing.T) {
	suite.Run(t, &DefaultFormDataDecoderImplTestSuite{})
}

func (t *DefaultFormDataDecoderImplTestSuite) SetupSuite() {
	t.decoder = &DefaultFormDataDecoderImpl{}
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecodeStringMap_Nil() {
	stringMap := t.decoder.decodeStringMap(nil)

	t.Equal(map[string]string{}, stringMap)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecodeStringMap_Empty() {
	stringMap := t.decoder.decodeStringMap(url.Values{})

	t.Equal(map[string]string{}, stringMap)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecodeStringMap_Values() {
	stringMap := t.decoder.decodeStringMap(url.Values{
		"first":  []string{"11", "12"},
		"second": []string{"21"},
	})

	t.Equal(map[string]string{
		"first":  "11 12",
		"second": "21",
	}, stringMap)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecodeUnknownInterface_Nil() {
	formData := formDataDecoderTestData{
		Text:   "some text",
		Number: 23,
		Slice:  []float64{1.0, 2.0},
	}

	result, err := t.decoder.decodeUnknownInterface(nil, formData)

	t.NoError(err)
	t.Equal(formDataDecoderTestData{}, result)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecodeUnknownInterface_Empty() {
	formData := formDataDecoderTestData{
		Text:   "some text",
		Number: 23,
		Slice:  []float64{1.0, 2.0},
	}

	result, err := t.decoder.decodeUnknownInterface(url.Values{}, formData)

	t.NoError(err)
	t.Equal(formDataDecoderTestData{}, result)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecodeUnknownInterface_Full() {
	formData := formDataDecoderTestData{
		Text:   "some text",
		Number: 23,
		Slice:  []float64{1.0, 2.0},
	}

	result, err := t.decoder.decodeUnknownInterface(url.Values{
		"text":   []string{" new text "},
		"number": []string{"10"},
	}, formData)

	t.NoError(err)
	t.Equal(formDataDecoderTestData{
		Text:   "new text",
		Number: 10,
	}, result)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecodeUnknownInterface_FullWithPointer() {
	formData := formDataDecoderTestData{
		Text:   "some text",
		Number: 23,
		Slice:  []float64{1.0, 2.0},
	}

	result, err := t.decoder.decodeUnknownInterface(url.Values{
		"text":   []string{" new text "},
		"number": []string{"10"},
	}, &formData)

	t.NoError(err)
	t.Equal(formDataDecoderTestData{
		Text:   "new text",
		Number: 10,
	}, result)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecode_StringMap() {
	stringMap, err := t.decoder.Decode(nil, nil, url.Values{
		"first":  []string{"11", "12"},
		"second": []string{"21"},
	}, map[string]string{})

	t.NoError(err)
	t.Equal(map[string]string{
		"first":  "11 12",
		"second": "21",
	}, stringMap)
}

func (t *DefaultFormDataDecoderImplTestSuite) TestDecode_UnknownInterface() {
	formData := formDataDecoderTestData{
		Text:   "some text",
		Number: 23,
		Slice:  []float64{1.0, 2.0},
	}

	result, err := t.decoder.Decode(nil, nil, url.Values{
		"text":   []string{" new text "},
		"number": []string{"10"},
	}, formData)

	t.NoError(err)
	t.Equal(formDataDecoderTestData{
		Text:   "new text",
		Number: 10,
	}, result)
}
