package application

import (
	"testing"

	"flamingo.me/flamingo/v3/core/form2/domain"
	"flamingo.me/flamingo/v3/core/form2/domain/mocks"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/suite"
)

type (
	FormHandlerBuilderImplTestSuite struct {
		suite.Suite

		builder *formHandlerBuilderImpl

		firstNamedService  *mocks.CompleteFormService
		secondNamedService *mocks.CompleteFormService
		service            *mocks.CompleteFormService

		firstNamedProvider  *mocks.FormDataProvider
		secondNamedProvider *mocks.FormDataProvider
		defaultProvider     *mocks.DefaultFormDataProvider
		provider            *mocks.FormDataProvider

		firstNamedDecoder  *mocks.FormDataDecoder
		secondNamedDecoder *mocks.FormDataDecoder
		defaultDecoder     *mocks.DefaultFormDataDecoder
		decoder            *mocks.FormDataDecoder

		firstNamedValidator  *mocks.FormDataValidator
		secondNamedValidator *mocks.FormDataValidator
		defaultValidator     *mocks.DefaultFormDataValidator
		validator            *mocks.FormDataValidator

		firstNamedExtension  *mocks.CompleteFormService
		secondNamedExtension *mocks.CompleteFormService

		validatorProvider *mocks.ValidatorProvider

		logger *flamingo.NullLogger
	}
)

func TestFormHandlerBuilderImplTestSuite(t *testing.T) {
	suite.Run(t, &FormHandlerBuilderImplTestSuite{})
}

func (t *FormHandlerBuilderImplTestSuite) SetupTest() {
	t.firstNamedService = &mocks.CompleteFormService{}
	t.secondNamedService = &mocks.CompleteFormService{}
	t.service = &mocks.CompleteFormService{}

	t.firstNamedProvider = &mocks.FormDataProvider{}
	t.secondNamedProvider = &mocks.FormDataProvider{}
	t.defaultProvider = &mocks.DefaultFormDataProvider{}
	t.provider = &mocks.FormDataProvider{}

	t.firstNamedDecoder = &mocks.FormDataDecoder{}
	t.secondNamedDecoder = &mocks.FormDataDecoder{}
	t.defaultDecoder = &mocks.DefaultFormDataDecoder{}
	t.decoder = &mocks.FormDataDecoder{}

	t.firstNamedValidator = &mocks.FormDataValidator{}
	t.secondNamedValidator = &mocks.FormDataValidator{}
	t.defaultValidator = &mocks.DefaultFormDataValidator{}
	t.validator = &mocks.FormDataValidator{}

	t.firstNamedExtension = &mocks.CompleteFormService{}
	t.secondNamedExtension = &mocks.CompleteFormService{}

	t.validatorProvider = &mocks.ValidatorProvider{}

	t.logger = &flamingo.NullLogger{}

	t.builder = &formHandlerBuilderImpl{
		namedFormServices: map[string]domain.FormService{
			"first":  t.firstNamedService,
			"second": t.secondNamedService,
		},
		namedFormDataProviders: map[string]domain.FormDataProvider{
			"first":  t.firstNamedProvider,
			"second": t.secondNamedProvider,
		},
		namedFormDataDecoders: map[string]domain.FormDataDecoder{
			"first":  t.firstNamedDecoder,
			"second": t.secondNamedDecoder,
		},
		namedFormDataValidators: map[string]domain.FormDataValidator{
			"first":  t.firstNamedValidator,
			"second": t.secondNamedValidator,
		},
		namedFormExtensions: map[string]domain.FormExtension{
			"first":  t.firstNamedExtension,
			"second": t.secondNamedExtension,
		},
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		validatorProvider:        t.validatorProvider,
		logger:                   t.logger,
	}
}

func (t *FormHandlerBuilderImplTestSuite) TearDownTest() {
	t.firstNamedService.AssertExpectations(t.T())
	t.secondNamedService.AssertExpectations(t.T())
	t.service.AssertExpectations(t.T())

	t.firstNamedProvider.AssertExpectations(t.T())
	t.secondNamedProvider.AssertExpectations(t.T())
	t.defaultProvider.AssertExpectations(t.T())
	t.provider.AssertExpectations(t.T())

	t.firstNamedDecoder.AssertExpectations(t.T())
	t.secondNamedDecoder.AssertExpectations(t.T())
	t.defaultDecoder.AssertExpectations(t.T())
	t.decoder.AssertExpectations(t.T())

	t.firstNamedValidator.AssertExpectations(t.T())
	t.secondNamedValidator.AssertExpectations(t.T())
	t.defaultValidator.AssertExpectations(t.T())
	t.validator.AssertExpectations(t.T())

	t.firstNamedExtension.AssertExpectations(t.T())
	t.secondNamedExtension.AssertExpectations(t.T())

	t.validatorProvider.AssertExpectations(t.T())
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormService_Panic() {
	t.Panics(func() {
		t.builder.Must(t.builder.SetFormService(nil))
	})
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormService_FormDataProvider() {
	t.Nil(t.builder.formDataProvider)
	t.Nil(t.builder.formDataDecoder)
	t.Nil(t.builder.formDataValidator)

	err := t.builder.SetFormService(t.provider)
	t.NoError(err)

	t.Exactly(t.provider, t.builder.formDataProvider)
	t.Nil(t.builder.formDataDecoder)
	t.Nil(t.builder.formDataValidator)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormService_FormDataDecoder() {
	t.Nil(t.builder.formDataProvider)
	t.Nil(t.builder.formDataDecoder)
	t.Nil(t.builder.formDataValidator)

	err := t.builder.SetFormService(t.decoder)
	t.NoError(err)

	t.Nil(t.builder.formDataProvider)
	t.Exactly(t.decoder, t.builder.formDataDecoder)
	t.Nil(t.builder.formDataValidator)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormService_FormDataValidator() {
	t.Nil(t.builder.formDataProvider)
	t.Nil(t.builder.formDataDecoder)
	t.Nil(t.builder.formDataValidator)

	err := t.builder.SetFormService(t.validator)
	t.NoError(err)

	t.Nil(t.builder.formDataProvider)
	t.Nil(t.builder.formDataDecoder)
	t.Exactly(t.validator, t.builder.formDataValidator)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormService_CompleteFormService() {
	t.Nil(t.builder.formDataProvider)
	t.Nil(t.builder.formDataDecoder)
	t.Nil(t.builder.formDataValidator)

	err := t.builder.SetFormService(t.service)
	t.NoError(err)

	t.Exactly(t.service, t.builder.formDataProvider)
	t.Exactly(t.service, t.builder.formDataDecoder)
	t.Exactly(t.service, t.builder.formDataValidator)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormService_Panic() {
	t.Panics(func() {
		t.builder.Must(t.builder.SetNamedFormService("third"))
	})
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormService_CompleteFormService() {
	t.Nil(t.builder.formDataProvider)
	t.Nil(t.builder.formDataDecoder)
	t.Nil(t.builder.formDataValidator)

	err := t.builder.SetNamedFormService("first")
	t.NoError(err)

	t.Exactly(t.firstNamedService, t.builder.formDataProvider)
	t.Exactly(t.firstNamedService, t.builder.formDataDecoder)
	t.Exactly(t.firstNamedService, t.builder.formDataValidator)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormDataProvider() {
	t.Nil(t.builder.formDataProvider)

	t.builder.SetFormDataProvider(t.provider)

	t.Exactly(t.provider, t.builder.formDataProvider)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormDataProvider_Panic() {
	t.Panics(func() {
		t.builder.Must(t.builder.SetNamedFormDataProvider("third"))
	})
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormDataProvider_Success() {
	t.Nil(t.builder.formDataProvider)

	err := t.builder.SetNamedFormDataProvider("first")
	t.NoError(err)

	t.Exactly(t.firstNamedProvider, t.builder.formDataProvider)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormDataDecoder() {
	t.Nil(t.builder.formDataDecoder)

	t.builder.SetFormDataDecoder(t.decoder)

	t.Exactly(t.decoder, t.builder.formDataDecoder)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormDataDecoder_Panic() {
	t.Panics(func() {
		t.builder.Must(t.builder.SetNamedFormDataDecoder("third"))
	})
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormDataDecoder_Success() {
	t.Nil(t.builder.formDataDecoder)

	err := t.builder.SetNamedFormDataDecoder("first")
	t.NoError(err)

	t.Exactly(t.firstNamedDecoder, t.builder.formDataDecoder)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetFormDataValidator() {
	t.Nil(t.builder.formDataValidator)

	t.builder.SetFormDataValidator(t.validator)

	t.Exactly(t.validator, t.builder.formDataValidator)
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormDataValidator_Panic() {
	t.Panics(func() {
		t.builder.Must(t.builder.SetNamedFormDataValidator("third"))
	})
}

func (t *FormHandlerBuilderImplTestSuite) TestSetNamedFormDataValidator_Success() {
	t.Nil(t.builder.formDataValidator)

	err := t.builder.SetNamedFormDataValidator("first")
	t.NoError(err)

	t.Exactly(t.firstNamedValidator, t.builder.formDataValidator)
}

func (t *FormHandlerBuilderImplTestSuite) TestAddFormExtension_Panic() {
	t.Panics(func() {
		t.builder.Must(t.builder.AddFormExtension(nil))
	})

	t.Panics(func() {
		t.builder.Must(t.builder.AddFormExtension("something wrong"))
	})
}

func (t *FormHandlerBuilderImplTestSuite) TestAddFormExtension_CompleteFormService() {
	t.Empty(t.builder.formExtensions)

	err := t.builder.AddFormExtension(t.service)
	t.NoError(err)

	t.Equal(map[string]domain.FormExtension{
		"CompleteFormService": t.service,
	}, t.builder.formExtensions)
}

func (t *FormHandlerBuilderImplTestSuite) TestAddNamedFormExtension_Panic() {
	t.Panics(func() {
		t.builder.Must(t.builder.AddNamedFormExtension("third"))
	})
}

func (t *FormHandlerBuilderImplTestSuite) TestAddNamedFormExtension_CompleteFormService() {
	t.Empty(t.builder.formExtensions)

	t.builder.AddNamedFormExtension("first")

	t.Equal(map[string]domain.FormExtension{
		"first": t.firstNamedExtension,
	}, t.builder.formExtensions)
}

func (t *FormHandlerBuilderImplTestSuite) TestBuild_Empty() {
	t.Equal(&formHandlerImpl{
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		formExtensions:           map[string]domain.FormExtension(nil),
		validatorProvider:        t.validatorProvider,
		logger:                   t.logger,
	}, t.builder.Build())
}

func (t *FormHandlerBuilderImplTestSuite) TestBuild_Full() {
	t.builder.SetFormDataProvider(t.provider)
	t.builder.SetFormDataDecoder(t.decoder)
	t.builder.SetFormDataValidator(t.validator)
	t.builder.AddFormExtension(t.service)

	t.Equal(&formHandlerImpl{
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		formDataProvider:         t.provider,
		formDataDecoder:          t.decoder,
		formDataValidator:        t.validator,
		formExtensions: map[string]domain.FormExtension{
			"CompleteFormService": t.service,
		},
		validatorProvider: t.validatorProvider,
		logger:            t.logger,
	}, t.builder.Build())
}
