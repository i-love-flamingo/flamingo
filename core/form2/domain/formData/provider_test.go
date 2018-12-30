package formData

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type (
	DefaultFormDataProviderImplTestSuite struct {
		suite.Suite

		provider *DefaultFormDataProviderImpl
	}
)

func TestDefaultFormDataProviderImplTestSuite(t *testing.T) {
	suite.Run(t, &DefaultFormDataProviderImplTestSuite{})
}

func (t *DefaultFormDataProviderImplTestSuite) SetupSuite() {
	t.provider = &DefaultFormDataProviderImpl{}
}

func (t *DefaultFormDataProviderImplTestSuite) TestGetFormData() {
	stringMap, err := t.provider.GetFormData(nil, nil)

	t.NoError(err)
	t.Equal(map[string]string{}, stringMap)
}
