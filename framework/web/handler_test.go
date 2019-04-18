package web

import (
	"context"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/zemirco/memorystore"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type (
	HandlerTestSuite struct {
		suite.Suite

		handler  *handler
		action   Action
		error    Action
		recorder *httptest.ResponseRecorder
	}
)

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, &HandlerTestSuite{})
}

func (t *HandlerTestSuite) SetupSuite() {
	t.action = func(ctx context.Context, req *Request) Result {
		return &Response{
			Status: http.StatusCreated,
		}
	}
	t.error = func(ctx context.Context, req *Request) Result {
		return &Response{
			Status: http.StatusNotFound,
		}
	}
}

func (t *HandlerTestSuite) SetupTest() {
	registry := NewRegistry()
	registry.HandleGet("test", t.action)
	registry.HandleAny(FlamingoNotfound, t.error)
	_, err := registry.Route("/test", "test")
	t.NoError(err)

	t.handler = &handler{
		routerRegistry: registry,
		eventRouter:    &flamingo.DefaultEventRouter{},
		logger:         &flamingo.NullLogger{},
		sessionStore:   memorystore.NewMemoryStore([]byte{}),
		sessionName:    "session",
	}

	t.recorder = httptest.NewRecorder()
}

func (t *HandlerTestSuite) TearDownTest() {
	t.handler = nil
}

func (t *HandlerTestSuite) TestServeHTTP_Found() {
	request, err := http.NewRequest(http.MethodGet, "/test", nil)
	t.NoError(err)

	t.handler.ServeHTTP(t.recorder, request)
	t.Equal(http.StatusCreated, t.recorder.Code)
}

func (t *HandlerTestSuite) TestServeHTTP_NotFound() {
	request, err := http.NewRequest(http.MethodGet, "/wrong", nil)
	t.NoError(err)

	t.handler.ServeHTTP(t.recorder, request)
	t.Equal(http.StatusNotFound, t.recorder.Code)
}
