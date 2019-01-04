package application

import (
	"testing"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/core/form2/domain/mocks"
	"flamingo.me/flamingo/framework/flamingo"
	"github.com/stretchr/testify/suite"
)

type (
	FormHandlerFactoryImplTestSuite struct {
		suite.Suite

		factory *FormHandlerFactoryImpl

		firstNamedService  *mocks.CompleteFormServiceWithName
		secondNamedService *mocks.CompleteFormServiceWithName
		service            *mocks.CompleteFormService

		firstNamedProvider  *mocks.FormDataProviderWithName
		secondNamedProvider *mocks.FormDataProviderWithName
		defaultProvider     *mocks.DefaultFormDataProvider
		provider            *mocks.FormDataProvider

		firstNamedDecoder  *mocks.FormDataDecoderWithName
		secondNamedDecoder *mocks.FormDataDecoderWithName
		defaultDecoder     *mocks.DefaultFormDataDecoder
		decoder            *mocks.FormDataDecoder

		firstNamedValidator  *mocks.FormDataValidatorWithName
		secondNamedValidator *mocks.FormDataValidatorWithName
		defaultValidator     *mocks.DefaultFormDataValidator
		validator            *mocks.FormDataValidator

		firstNamedExtension  *mocks.CompleteFormServiceWithName
		secondNamedExtension *mocks.CompleteFormServiceWithName

		validatorProvider *mocks.ValidatorProvider

		logger *flamingo.NullLogger
	}
)

func TestFormHandlerFactoryImplTestSuite(t *testing.T) {
	suite.Run(t, &FormHandlerFactoryImplTestSuite{})
}

func (t *FormHandlerFactoryImplTestSuite) SetupTest() {
	t.firstNamedService = &mocks.CompleteFormServiceWithName{}
	t.secondNamedService = &mocks.CompleteFormServiceWithName{}
	t.service = &mocks.CompleteFormService{}

	t.firstNamedProvider = &mocks.FormDataProviderWithName{}
	t.secondNamedProvider = &mocks.FormDataProviderWithName{}
	t.defaultProvider = &mocks.DefaultFormDataProvider{}
	t.provider = &mocks.FormDataProvider{}

	t.firstNamedDecoder = &mocks.FormDataDecoderWithName{}
	t.secondNamedDecoder = &mocks.FormDataDecoderWithName{}
	t.defaultDecoder = &mocks.DefaultFormDataDecoder{}
	t.decoder = &mocks.FormDataDecoder{}

	t.firstNamedValidator = &mocks.FormDataValidatorWithName{}
	t.secondNamedValidator = &mocks.FormDataValidatorWithName{}
	t.defaultValidator = &mocks.DefaultFormDataValidator{}
	t.validator = &mocks.FormDataValidator{}

	t.firstNamedExtension = &mocks.CompleteFormServiceWithName{}
	t.secondNamedExtension = &mocks.CompleteFormServiceWithName{}

	t.validatorProvider = &mocks.ValidatorProvider{}

	t.logger = &flamingo.NullLogger{}

	t.factory = &FormHandlerFactoryImpl{}
	t.factory.Inject(
		[]domain.FormServiceWithName{
			t.firstNamedService,
			t.secondNamedService,
		},
		[]domain.FormDataProviderWithName{
			t.firstNamedProvider,
			t.secondNamedProvider,
		},
		[]domain.FormDataDecoderWithName{
			t.firstNamedDecoder,
			t.secondNamedDecoder,
		},
		[]domain.FormDataValidatorWithName{
			t.firstNamedValidator,
			t.secondNamedValidator,
		},
		[]domain.FormExtensionWithName{
			t.firstNamedExtension,
			t.secondNamedExtension,
		},
		t.defaultProvider,
		t.defaultDecoder,
		t.defaultValidator,
		t.validatorProvider,
		t.logger,
	)
}

func (t *FormHandlerFactoryImplTestSuite) TearDownTest() {
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

func (t *FormHandlerFactoryImplTestSuite) TestCreateSimpleFormHandler() {
	t.Equal(&formHandlerImpl{
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		formExtensions:           []interface{}(nil),
		validatorProvider:        t.validatorProvider,
		logger:                   t.logger,
	}, t.factory.CreateSimpleFormHandler())
}

func (t *FormHandlerFactoryImplTestSuite) TestCreateFormHandlerWithFormService() {
	t.Equal(&formHandlerImpl{
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		formDataProvider:         t.service,
		formDataDecoder:          t.service,
		formDataValidator:        t.service,
		formExtensions: []interface{}{
			t.validator,
			t.secondNamedValidator,
		},
		validatorProvider: t.validatorProvider,
		logger:            t.logger,
	}, t.factory.CreateFormHandlerWithFormService(t.service, t.validator, t.secondNamedValidator))
}

func (t *FormHandlerFactoryImplTestSuite) TestCreateFormHandlerWithFormServices() {
	t.Equal(&formHandlerImpl{
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		formDataProvider:         t.provider,
		formDataDecoder:          t.decoder,
		formDataValidator:        t.validator,
		formExtensions: []interface{}{
			t.service,
			t.firstNamedService,
		},
		validatorProvider: t.validatorProvider,
		logger:            t.logger,
	}, t.factory.CreateFormHandlerWithFormServices(t.provider, t.decoder, t.validator, t.service, t.firstNamedService))
}

func (t *FormHandlerFactoryImplTestSuite) TestGetFormHandlerBuilder() {
	t.Equal(&formHandlerBuilderImpl{
		namedFormServices: []domain.FormServiceWithName{
			t.firstNamedService,
			t.secondNamedService,
		},
		namedFormDataProviders: []domain.FormDataProviderWithName{
			t.firstNamedProvider,
			t.secondNamedProvider,
		},
		namedFormDataDecoders: []domain.FormDataDecoderWithName{
			t.firstNamedDecoder,
			t.secondNamedDecoder,
		},
		namedFormDataValidators: []domain.FormDataValidatorWithName{
			t.firstNamedValidator,
			t.secondNamedValidator,
		},
		namedFormExtensions: []domain.FormExtensionWithName{
			t.firstNamedExtension,
			t.secondNamedExtension,
		},
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		validatorProvider:        t.validatorProvider,
		logger:                   t.logger,
	}, t.factory.GetFormHandlerBuilder())
}
