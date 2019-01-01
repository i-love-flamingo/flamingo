package application

import (
	"context"
	"errors"
	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"

	"flamingo.me/flamingo/core/form2/domain/mocks"
	"flamingo.me/flamingo/framework/web"
)

type (
	FormHandlerImplTestSuite struct {
		suite.Suite

		handler *formHandlerImpl

		provider          *mocks.FormDataProvider
		decoder           *mocks.FormDataDecoder
		validator         *mocks.FormDataValidator
		firstExtension    *mocks.CompleteFormService
		secondExtension   *mocks.FormDataProvider
		thirdExtension    *mocks.FormDataDecoder
		fourthExtension   *mocks.FormDataValidator
		validatorProvider *mocks.ValidatorProvider
		logger            *flamingo.NullLogger

		context context.Context
		request *web.Request
	}
)

func TestFormHandlerImplTestSuite(t *testing.T) {
	suite.Run(t, &FormHandlerImplTestSuite{})
}

func (t *FormHandlerImplTestSuite) SetupSuite() {
	t.context = context.Background()
}

func (t *FormHandlerImplTestSuite) SetupTest() {
	t.provider = &mocks.FormDataProvider{}
	t.decoder = &mocks.FormDataDecoder{}
	t.validator = &mocks.FormDataValidator{}
	t.firstExtension = &mocks.CompleteFormService{}
	t.secondExtension = &mocks.FormDataProvider{}
	t.thirdExtension = &mocks.FormDataDecoder{}
	t.fourthExtension = &mocks.FormDataValidator{}
	t.validatorProvider = &mocks.ValidatorProvider{}
	t.logger = &flamingo.NullLogger{}

	t.handler = &formHandlerImpl{
		formDataProvider:  t.provider,
		formDataDecoder:   t.decoder,
		formDataValidator: t.validator,
		formExtensions: []interface{}{
			t.firstExtension,
			t.secondExtension,
			t.thirdExtension,
			t.fourthExtension,
		},
		validatorProvider: t.validatorProvider,
		logger:            t.logger,
	}

	t.request = web.RequestFromRequest(&http.Request{}, nil)
}

func (t *FormHandlerImplTestSuite) TearDownTest() {
	t.provider.AssertExpectations(t.T())
	t.decoder.AssertExpectations(t.T())
	t.validator.AssertExpectations(t.T())
	t.firstExtension.AssertExpectations(t.T())
	t.secondExtension.AssertExpectations(t.T())
	t.thirdExtension.AssertExpectations(t.T())
	t.fourthExtension.AssertExpectations(t.T())
	t.validatorProvider.AssertExpectations(t.T())
}

func (t *FormHandlerImplTestSuite) TestHandleUnsubmittedForm_Error() {
	t.provider.On("GetFormData", t.context, t.request).Return(nil, errors.New("error")).Once()

	result, err := t.handler.HandleUnsubmittedForm(t.context, t.request)
	t.Equal(domain.FormError("error"), err)
	t.Nil(result)
}

func (t *FormHandlerImplTestSuite) TestHandleUnsubmittedForm_Success() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]int{}, nil).Once()

	result, err := t.handler.HandleUnsubmittedForm(t.context, t.request)
	t.NoError(err)

	form := domain.NewForm(false, map[string][]domain.ValidationRule{})
	form.Data = map[string]int{}

	t.Equal(&form, result)
}

func (t *FormHandlerImplTestSuite) TestExtractValidationRules_NotStruct() {
	t.Equal(map[string][]domain.ValidationRule{}, t.handler.extractValidationRules(nil))
	t.Equal(map[string][]domain.ValidationRule{}, t.handler.extractValidationRules("string"))
	t.Equal(map[string][]domain.ValidationRule{}, t.handler.extractValidationRules(1))
	t.Equal(map[string][]domain.ValidationRule{}, t.handler.extractValidationRules(map[string]interface{}{}))
}

func (t *FormHandlerImplTestSuite) TestExtractValidationRules_Struct() {
	t.Equal(map[string][]domain.ValidationRule{
		"first": {
			{
				Name: "required",
			},
			{
				Name:  "gte",
				Value: "10",
			},
		},
		"second": {
			{
				Name:  "gte",
				Value: "10",
			},
		},
	}, t.handler.extractValidationRules(struct {
		First  string `form:"first" validate:"required,gte=10"`
		Second string `form:"second" validate:"omitempty,gte=10"`
		Third  string `form:"-" validate:"required,gte=10"`
		Fourth string `form:"fourth" validate:""`
		Fifth  string `form:"fifth"`
		Sixth  string `validate:"required,gte=10"`
	}{}))
}

func (t *FormHandlerImplTestSuite) TestGetPostValues_Error() {
	t.request.Request().Method = http.MethodPost

	values, err := t.handler.getPostValues(t.request)
	t.Error(err)
	t.Nil(values)
}

func (t *FormHandlerImplTestSuite) TestGetPostValues_Success() {
	t.request.Request().Method = http.MethodPost
	t.request.Request().PostForm = url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}

	values, err := t.handler.getPostValues(t.request)
	t.NoError(err)
	t.Equal(&url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, values)
}

func (t *FormHandlerImplTestSuite) TestProcessExtension_ProviderError() {
	t.secondExtension.On("GetFormData", t.context, t.request).Return(nil, errors.New("error")).Once()

	err := t.handler.processExtension(t.context, t.request, url.Values{}, t.secondExtension, &domain.Form{})
	t.Error(err)
}

func (t *FormHandlerImplTestSuite) TestProcessExtension_ProviderSuccess() {
	t.secondExtension.On("GetFormData", t.context, t.request).Return(map[string]int{}, nil).Once()

	err := t.handler.processExtension(t.context, t.request, url.Values{}, t.secondExtension, &domain.Form{})
	t.NoError(err)
}

func (t *FormHandlerImplTestSuite) TestProcessExtension_DecoderError() {
	t.thirdExtension.On("Decode", t.context, t.request, url.Values{}, nil).Return(nil, errors.New("error")).Once()

	err := t.handler.processExtension(t.context, t.request, url.Values{}, t.thirdExtension, &domain.Form{})
	t.Error(err)
}

func (t *FormHandlerImplTestSuite) TestProcessExtension_DecoderSuccess() {
	t.thirdExtension.On("Decode", t.context, t.request, url.Values{}, nil).Return(map[string]int{}, nil).Once()

	err := t.handler.processExtension(t.context, t.request, url.Values{}, t.thirdExtension, &domain.Form{})
	t.NoError(err)
}

func (t *FormHandlerImplTestSuite) TestProcessExtension_ValidatorError() {
	t.fourthExtension.On("Validate", t.context, t.request, t.validatorProvider, nil).Return(nil, errors.New("error")).Once()

	err := t.handler.processExtension(t.context, t.request, url.Values{}, t.fourthExtension, &domain.Form{})
	t.Error(err)
}

func (t *FormHandlerImplTestSuite) TestProcessExtension_ValidatorSuccess() {
	t.fourthExtension.On("Validate", t.context, t.request, t.validatorProvider, nil).Return(&domain.ValidationInfo{}, nil).Once()

	err := t.handler.processExtension(t.context, t.request, url.Values{}, t.fourthExtension, &domain.Form{})
	t.NoError(err)
}

func (t *FormHandlerImplTestSuite) TestHandleSubmittedForm_GetFormDataError() {
	t.provider.On("GetFormData", t.context, t.request).Return(nil, errors.New("error")).Once()

	result, err := t.handler.HandleSubmittedForm(t.context, t.request)
	t.Equal(domain.FormError("error"), err)
	t.Nil(result)
}

func (t *FormHandlerImplTestSuite) TestHandleSubmittedForm_PostValueProcessingError() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]string{}, nil).Once()

	t.request.Request().Method = http.MethodPost

	result, err := t.handler.HandleSubmittedForm(t.context, t.request)
	t.Equal(domain.FormError("missing form body"), err)
	t.Nil(result)
}

func (t *FormHandlerImplTestSuite) TestHandleSubmittedForm_DecodeError() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]string{}, nil).Once()

	t.request.Request().Method = http.MethodPost
	t.request.Request().PostForm = url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}

	t.decoder.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, map[string]string{}).Return(nil, errors.New("error")).Once()

	result, err := t.handler.HandleSubmittedForm(t.context, t.request)
	t.Equal(domain.FormError("error"), err)
	t.Nil(result)
}

func (t *FormHandlerImplTestSuite) TestHandleSubmittedForm_ValidateError() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]string{}, nil).Once()

	t.request.Request().Method = http.MethodPost
	t.request.Request().PostForm = url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}

	t.decoder.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, map[string]string{}).Return(map[string]string{
		"first":  "first",
		"second": "second",
	}, nil).Once()

	t.validator.On("Validate", t.context, t.request, t.validatorProvider, map[string]string{
		"first":  "first",
		"second": "second",
	}).Return(nil, errors.New("error")).Once()

	result, err := t.handler.HandleSubmittedForm(t.context, t.request)
	t.Equal(domain.FormError("error"), err)
	t.Nil(result)
}

func (t *FormHandlerImplTestSuite) TestHandleSubmittedForm_FormExtensionError() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]string{}, nil).Once()

	t.request.Request().Method = http.MethodPost
	t.request.Request().PostForm = url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}

	t.decoder.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, map[string]string{}).Return(map[string]string{
		"first":  "first",
		"second": "second",
	}, nil).Once()

	t.validator.On("Validate", t.context, t.request, t.validatorProvider, map[string]string{
		"first":  "first",
		"second": "second",
	}).Return(&domain.ValidationInfo{}, nil).Once()

	t.firstExtension.On("GetFormData", t.context, t.request).Return(nil, errors.New("error")).Once()

	result, err := t.handler.HandleSubmittedForm(t.context, t.request)
	t.Equal(domain.FormError("error"), err)
	t.Nil(result)
}

func (t *FormHandlerImplTestSuite) TestHandleSubmittedForm_Success() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]string{}, nil).Once()

	t.request.Request().Method = http.MethodPost
	t.request.Request().PostForm = url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}

	t.decoder.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, map[string]string{}).Return(map[string]string{
		"first":  "first",
		"second": "second",
	}, nil).Once()

	t.validator.On("Validate", t.context, t.request, t.validatorProvider, map[string]string{
		"first":  "first",
		"second": "second",
	}).Return(&domain.ValidationInfo{}, nil).Once()

	t.firstExtension.On("GetFormData", t.context, t.request).Return(map[string]float64{}, nil).Once()
	t.firstExtension.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, map[string]float64{}).Return(map[string]float64{}, nil).Once()
	t.firstExtension.On("Validate", t.context, t.request, t.validatorProvider, map[string]float64{}).Return(&domain.ValidationInfo{}, nil).Once()
	t.secondExtension.On("GetFormData", t.context, t.request).Return(map[string]int{}, nil).Once()
	t.thirdExtension.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, nil).Return(map[string]int{}, nil).Once()
	t.fourthExtension.On("Validate", t.context, t.request, t.validatorProvider, nil).Return(&domain.ValidationInfo{}, nil).Once()

	result, err := t.handler.HandleSubmittedForm(t.context, t.request)
	t.NoError(err)

	form := domain.NewForm(true, map[string][]domain.ValidationRule{})
	form.Data = map[string]string{
		"first":  "first",
		"second": "second",
	}

	t.Equal(&form, result)
}

func (t *FormHandlerImplTestSuite) TestHandleForm_Unsubmitted() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]int{}, nil).Once()

	result, err := t.handler.HandleUnsubmittedForm(t.context, t.request)
	t.NoError(err)

	form := domain.NewForm(false, map[string][]domain.ValidationRule{})
	form.Data = map[string]int{}

	t.Equal(&form, result)
}

func (t *FormHandlerImplTestSuite) TestHandleForm_Submitted() {
	t.provider.On("GetFormData", t.context, t.request).Return(map[string]string{}, nil).Once()

	t.request.Request().Method = http.MethodPost
	t.request.Request().PostForm = url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}

	t.decoder.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, map[string]string{}).Return(map[string]string{
		"first":  "first",
		"second": "second",
	}, nil).Once()

	t.validator.On("Validate", t.context, t.request, t.validatorProvider, map[string]string{
		"first":  "first",
		"second": "second",
	}).Return(&domain.ValidationInfo{}, nil).Once()

	t.firstExtension.On("GetFormData", t.context, t.request).Return(map[string]float64{}, nil).Once()
	t.firstExtension.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, map[string]float64{}).Return(map[string]float64{}, nil).Once()
	t.firstExtension.On("Validate", t.context, t.request, t.validatorProvider, map[string]float64{}).Return(&domain.ValidationInfo{}, nil).Once()
	t.secondExtension.On("GetFormData", t.context, t.request).Return(map[string]int{}, nil).Once()
	t.thirdExtension.On("Decode", t.context, t.request, url.Values{
		"first":  []string{"first"},
		"second": []string{"second"},
	}, nil).Return(map[string]int{}, nil).Once()
	t.fourthExtension.On("Validate", t.context, t.request, t.validatorProvider, nil).Return(&domain.ValidationInfo{}, nil).Once()

	result, err := t.handler.HandleForm(t.context, t.request)
	t.NoError(err)

	form := domain.NewForm(true, map[string][]domain.ValidationRule{})
	form.Data = map[string]string{
		"first":  "first",
		"second": "second",
	}

	t.Equal(&form, result)
}
