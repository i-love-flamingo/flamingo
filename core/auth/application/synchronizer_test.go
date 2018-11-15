package application

import (
	"testing"

	"errors"

	"flamingo.me/flamingo/core/auth/application/mocks"
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	SynchronizerTestSuite struct {
		suite.Suite

		synchronizerImpl *SynchronizerImpl

		store  *mocks.Store
		logger *flamingo.NullLogger

		session *sessions.Session
	}
)

func TestSynchronizerTestSuite(t *testing.T) {
	suite.Run(t, &SynchronizerTestSuite{})
}

func (t *SynchronizerTestSuite) SetupTest() {
	t.logger = &flamingo.NullLogger{}
	t.store = &mocks.Store{}
	t.synchronizerImpl = &SynchronizerImpl{}
	t.synchronizerImpl.Inject(t.store, t.logger)

	t.session = sessions.NewSession(nil, "-")
	t.session.ID = "someSessionID"
}

func (t *SynchronizerTestSuite) TearDownTest() {
	t.store.AssertExpectations(t.T())
	t.store = nil
	t.logger = nil
	t.synchronizerImpl = nil
}

func (t *SynchronizerTestSuite) TestInsert_DestroyError() {
	user := domain.User{
		Sub: "sub",
	}

	t.store.On("DestroySessionsForUser", user).Return(errors.New("error")).Once()
	t.store.On("SetHashAndSessionIdForUser", user, mock.Anything, "someSessionID").Return(nil).Once()

	t.NoError(t.synchronizerImpl.Insert(user, t.session))
}

func (t *SynchronizerTestSuite) TestInsert_Success() {
	user := domain.User{
		Sub: "sub",
	}

	t.store.On("DestroySessionsForUser", user).Return(nil).Once()
	t.store.On("SetHashAndSessionIdForUser", user, mock.Anything, "someSessionID").Return(nil).Once()

	t.NoError(t.synchronizerImpl.Insert(user, t.session))
}

func (t *SynchronizerTestSuite) TestInsert_SetError() {
	user := domain.User{
		Sub: "sub",
	}

	t.store.On("DestroySessionsForUser", user).Return(nil).Once()
	t.store.On("SetHashAndSessionIdForUser", user, mock.Anything, "someSessionID").Return(errors.New("error")).Once()

	t.Error(t.synchronizerImpl.Insert(user, t.session))
}

func (t *SynchronizerTestSuite) TestIsActive_Error() {
	user := domain.User{
		Sub: "sub",
	}

	t.store.On("GetHashByUser", user).Return("", errors.New("error")).Once()

	result, err := t.synchronizerImpl.IsActive(user, t.session)
	t.Error(err)
	t.False(result)
}

func (t *SynchronizerTestSuite) TestIsActive_NotActive() {
	t.session.Values[hashKey] = "hash"
	user := domain.User{
		Sub: "sub",
	}

	t.store.On("GetHashByUser", user).Return("different", nil).Once()

	result, err := t.synchronizerImpl.IsActive(user, t.session)
	t.NoError(err)
	t.False(result)
}

func (t *SynchronizerTestSuite) TestIsActive_Active() {
	t.session.Values[hashKey] = "hash"
	user := domain.User{
		Sub: "sub",
	}

	t.store.On("GetHashByUser", user).Return("hash", nil).Once()

	result, err := t.synchronizerImpl.IsActive(user, t.session)
	t.NoError(err)
	t.True(result)
}
