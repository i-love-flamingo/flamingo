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

func TestFormHandlerFactoryImplTestSuite(t *testing.T) {
	suite.Run(t, &FormHandlerFactoryImplTestSuite{})
}

func (t *FormHandlerFactoryImplTestSuite) SetupTest() {
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

	t.factory = &FormHandlerFactoryImpl{}
	t.factory.Inject(
		map[string]domain.FormService{
			"first":  t.firstNamedService,
			"second": t.secondNamedService,
		},
		map[string]domain.FormDataProvider{
			"first":  t.firstNamedProvider,
			"second": t.secondNamedProvider,
		},
		map[string]domain.FormDataDecoder{
			"first":  t.firstNamedDecoder,
			"second": t.secondNamedDecoder,
		},
		map[string]domain.FormDataValidator{
			"first":  t.firstNamedValidator,
			"second": t.secondNamedValidator,
		},
		map[string]domain.FormExtension{
			"first":  t.firstNamedExtension,
			"second": t.secondNamedExtension,
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
		formExtensions:           map[string]domain.FormExtension(nil),
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
		formExtensions: map[string]domain.FormExtension{
			"first":  t.firstNamedExtension,
			"second": t.secondNamedExtension,
		},
		validatorProvider: t.validatorProvider,
		logger:            t.logger,
	}, t.factory.CreateFormHandlerWithFormService(t.service, "first", "second"))
}

func (t *FormHandlerFactoryImplTestSuite) TestCreateFormHandlerWithFormServices() {
	t.Equal(&formHandlerImpl{
		defaultFormDataProvider:  t.defaultProvider,
		defaultFormDataDecoder:   t.defaultDecoder,
		defaultFormDataValidator: t.defaultValidator,
		formDataProvider:         t.provider,
		formDataDecoder:          t.decoder,
		formDataValidator:        t.validator,
		formExtensions: map[string]domain.FormExtension{
			"first":  t.firstNamedExtension,
			"second": t.secondNamedExtension,
		},
		validatorProvider: t.validatorProvider,
		logger:            t.logger,
	}, t.factory.CreateFormHandlerWithFormServices(t.provider, t.decoder, t.validator, "first", "second"))
}

func (t *FormHandlerFactoryImplTestSuite) TestGetFormHandlerBuilder() {
	t.Equal(&formHandlerBuilderImpl{
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
	}, t.factory.GetFormHandlerBuilder())
}
